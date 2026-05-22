require 'rails_helper'

RSpec.describe "Api::V0::BatchEducationEnrollments", type: :request do
  let(:coordinator_mock) { instance_double(Education::ServiceCoordinator) }
  let(:batch_id) { "test-batch-123" }

  before do
    allow(Education::ServiceCoordinator).to receive(:new).and_return(coordinator_mock)
    # Mock Reporter to avoid actual reporting during tests
    allow(Reporting::Reporter).to receive(:new).and_return(double(report: nil))
  end

  describe "GET /api/v0/batch-education-enrollments/:id" do
    it "returns the batch status" do
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

      get "/api/v0/batch-education-enrollments/#{batch_id}"

      expect(response).to have_http_status(:ok)
      expect(JSON.parse(response.body)["batch_job_id"]).to eq(batch_id)
    end

    it "returns 404 when batch is not found" do
      expect(coordinator_mock).to receive(:get_batch_status).with("invalid").and_raise(ActiveRecord::RecordNotFound)

      get "/api/v0/batch-education-enrollments/invalid"

      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)["error"]).to eq("Batch not found")
    end
  end

  describe "GET /api/v0/batch-education-enrollments/:batchJobId/details" do
    it "returns the batch details" do
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

      get "/api/v0/batch-education-enrollments/#{batch_id}/details"

      expect(response).to have_http_status(:ok)
      expect(JSON.parse(response.body)["batch_job_id"]).to eq(batch_id)
      expect(JSON.parse(response.body)["results"]).to be_an(Array)
    end

    it "returns 404 when batch is not found" do
      expect(coordinator_mock).to receive(:get_batch_details).with("invalid").and_raise(ActiveRecord::RecordNotFound)

      get "/api/v0/batch-education-enrollments/invalid/details"

      expect(response).to have_http_status(:not_found)
      expect(JSON.parse(response.body)["error"]).to eq("Batch not found")
    end
  end
end
