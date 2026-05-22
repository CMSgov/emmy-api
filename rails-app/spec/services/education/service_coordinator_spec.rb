require 'rails_helper'

module Education
  RSpec.describe ServiceCoordinator do
    let(:coordinator) { ServiceCoordinator.new }
    let(:batch_id) { 'test-batch-123' }
    let(:params) do
      {
        batchId: batch_id,
        submittedBy: 'test-user',
        callbackUrl: 'http://example.com/callback',
        students: [
          {
            recordId: 'rec1',
            firstName: 'John',
            lastName: 'Doe',
            dateOfBirth: '1990-01-01',
            ssn: '123456789'
          }
        ]
      }
    end

    let(:batch_student_request) { Education::BatchStudentRequest.new(params) }

    let(:sqs_client) { instance_double(Aws::SQS::Client) }
    let(:queue_url) { 'http://sqs.test/batch-queue' }

    before do
      # Mock encryption config to avoid errors if ENV is not set in test environment
      allow(ActiveRecord::Encryption.config).to receive(:primary_key).and_return("test-primary-key")
      allow(ActiveRecord::Encryption.config).to receive(:deterministic_key).and_return("test-deterministic-key")
      allow(ActiveRecord::Encryption.config).to receive(:key_derivation_salt).and_return("test-salt")

      allow(ENV).to receive(:[]).and_call_original
      allow(ENV).to receive(:[]).with("BATCH_SQS_QUEUE_URL").and_return(queue_url)
      allow(Aws::SQS::Client).to receive(:new).and_return(sqs_client)
    end

    describe '#register_batch' do
      it 'creates a batch and its students' do
        expect(sqs_client).to receive(:send_message).with(hash_including(queue_url: queue_url))

        expect {
          coordinator.register_batch(batch_student_request)
        }.to change(Education::EnrollmentBatch, :count).by(1)
         .and change(Education::BatchStudent, :count).by(1)
      end

      it 'enqueues an SQS message with the enrollment_batch_id' do
        batch = nil
        expect(sqs_client).to receive(:send_message) do |args|
          expect(args[:queue_url]).to eq(queue_url)
          body = JSON.parse(args[:message_body])
          expect(body['enrollment_batch_id']).to eq(batch.id)
        end

        batch = coordinator.register_batch(batch_student_request)
      end

      it 'handles SQS errors gracefully' do
        allow(sqs_client).to receive(:send_message).and_raise(StandardError.new("SQS failure"))
        expect(Rails.logger).to receive(:error).with(/Failed to enqueue batch job to SQS: SQS failure/)

        expect {
          coordinator.register_batch(batch_student_request)
        }.not_to raise_error
      end
    end

    describe '#get_batch_status' do
      let!(:batch) { Education::EnrollmentBatch.create!(batch_id: batch_id, submitted_by: 'test-user', callback_url: 'http://example.com', status: :processing) }
      let!(:student) { Education::BatchStudent.create!(education_enrollment_batch: batch, record_id: 'rec1', status: :success, first_name: 'J', last_name: 'D', date_of_birth: '1990-01-01', ssn: '123') }

      it 'returns the status of a batch' do
        status = coordinator.get_batch_status(batch_id)
        expect(status[:batch_job_id]).to eq(batch_id)
        expect(status[:status]).to eq('PROCESSING')
        expect(status[:total_records]).to eq(1)
        expect(status[:success_count]).to eq(1)
      end

      it 'raises RecordNotFound if batch does not exist' do
        expect {
          coordinator.get_batch_status('non-existent')
        }.to raise_error(ActiveRecord::RecordNotFound)
      end
    end

    describe '#get_batch_details' do
      let!(:batch) { Education::EnrollmentBatch.create!(batch_id: batch_id, submitted_by: 'test-user', callback_url: 'http://example.com', status: :completed) }
      let!(:student) { Education::BatchStudent.create!(education_enrollment_batch: batch, record_id: 'rec1', status: :success, first_name: 'J', last_name: 'D', date_of_birth: '1990-01-01', ssn: '123') }
      let!(:result) { Education::BatchStudentResult.create!(education_batch_student: student, found_enrollment: true, results: [ { school: 'Test' } ]) }

      it 'returns the details of a batch' do
        details = coordinator.get_batch_details(batch_id)
        expect(details[:batch_job_id]).to eq(batch_id)
        expect(details[:results].first[:record_id]).to eq('rec1')
        expect(details[:results].first[:found_enrollment]).to be(true)
      end
    end
  end
end
