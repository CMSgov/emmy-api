require 'swagger_helper'

RSpec.describe "Api::V0::BatchEducationEnrollments", type: :request do
  let(:coordinator_mock) { instance_double(Education::ServiceCoordinator) }
  let(:batch_id) { "test-batch-123" }

  before do
    allow(Education::ServiceCoordinator).to receive(:new).and_return(coordinator_mock)
    # Mock Reporter to avoid actual reporting during tests
    allow(Reporting::Reporter).to receive(:new).and_return(double(report: nil))
  end

  path '/api/v0/batch-education-enrollments' do
    post 'Register a new batch of education enrollments' do
      tags 'Batch Education Enrollments'
      consumes 'application/json'
      produces 'application/json'
      security [ oauth2: [] ]
      let(:Authorization) { 'Bearer <token>' }
      parameter name: :batch_params, in: :body, schema: Education::BatchStudentRequest.to_swagger_schema

      response '201', 'batch registration successful' do
        schema Education::BatchStudentCreatedResponse.to_swagger_schema

        let(:batch_params) do
          {
            batchId: batch_id,
            submittedBy: 'user@example.com',
            callbackUrl: 'https://example.com/callback',
            students: [
              {
                recordId: 'rec1',
                firstName: 'John',
                lastName: 'Doe',
                dateOfBirth: '1990-01-01'
              }
            ]
          }
        end

        before do
          expect(coordinator_mock).to receive(:register_batch).and_return(double(id: 'uuid-123'))
        end

        after do |example|
          example.metadata[:response][:content] = {
            'application/json' => {
              example: JSON.parse(response.body, symbolize_names: true)
            }
          }
        end
        run_test!
      end

      response '400', 'missing required fields' do
        let(:batch_params) { { batchId: batch_id } }
        run_test!
      end
    end
  end

  path '/api/v0/batch-education-enrollments/{id}' do
    parameter name: :id, in: :path, type: :string

    get 'Get the status of a batch' do
      tags 'Batch Education Enrollments'
      produces 'application/json'
      security [ oauth2: [] ]
      let(:Authorization) { 'Bearer <token>' }

      response '200', 'successful' do
        schema Education::BatchStatusResponse.to_swagger_schema

        let(:id) { batch_id }

        before do
          status_data = {
            batch_job_id: batch_id,
            status: "PROCESSING",
            submitted_at: Time.now,
            updated_at: Time.now,
            total_records: 1,
            processed_records: 0,
            success_count: 0,
            failure_count: 0
          }
          expect(coordinator_mock).to receive(:get_batch_status).with(batch_id).and_return(status_data)
        end

        after do |example|
          example.metadata[:response][:content] = {
            'application/json' => {
              example: JSON.parse(response.body, symbolize_names: true)
            }
          }
        end
        run_test!
      end

      response '404', 'batch not found' do
        let(:id) { 'invalid' }
        before do
          expect(coordinator_mock).to receive(:get_batch_status).with('invalid').and_raise(ActiveRecord::RecordNotFound)
        end
        run_test!
      end
    end
  end

  path '/api/v0/batch-education-enrollments/{batchJobId}/details' do
    parameter name: :batchJobId, in: :path, type: :string

    get 'Get the details of a batch' do
      tags 'Batch Education Enrollments'
      produces 'application/json'
      security [ oauth2: [] ]
      let(:Authorization) { 'Bearer <token>' }

      response '200', 'successful' do
        schema Education::BatchDetailsResponse.to_swagger_schema

        let(:batchJobId) { batch_id }

        before do
          details_data = {
            batch_job_id: batch_id,
            results: [
              {
                record_id: "rec1",
                status: "SUCCESS",
                found_enrollment: true,
                results: []
              }
            ]
          }
          expect(coordinator_mock).to receive(:get_batch_details).with(batch_id).and_return(details_data)
        end

        after do |example|
          example.metadata[:response][:content] = {
            'application/json' => {
              example: JSON.parse(response.body, symbolize_names: true)
            }
          }
        end
        run_test!
      end

      response '404', 'batch not found' do
        let(:batchJobId) { 'invalid' }
        before do
          expect(coordinator_mock).to receive(:get_batch_details).with('invalid').and_raise(ActiveRecord::RecordNotFound)
        end
        run_test!
      end
    end
  end
end
