module Education
  class BatchStatusResponse
    attr_accessor :batch_job_id, :status, :submitted_at, :updated_at, :total_records, :processed_records, :success_count, :failure_count

    def initialize(params = {})
      @batch_job_id = params[:batch_job_id]
      @status = params[:status]
      @submitted_at = params[:submitted_at]
      @updated_at = params[:updated_at]
      @total_records = params[:total_records]
      @processed_records = params[:processed_records]
      @success_count = params[:success_count]
      @failure_count = params[:failure_count]
    end

    def as_json(options = {})
      {
        batch_job_id: batch_job_id,
        status: status,
        submitted_at: submitted_at,
        updated_at: updated_at,
        total_records: total_records,
        processed_records: processed_records,
        success_count: success_count,
        failure_count: failure_count
      }
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          batch_job_id: { type: :string, example: 'batch-2023-05-22-001' },
          status: {
            type: :string,
            enum: Education::EnrollmentBatch.statuses.values,
            example: 'PROCESSING'
          },
          submitted_at: { type: :string, format: 'date-time', example: '2023-05-22T08:53:00Z' },
          updated_at: { type: :string, format: 'date-time', example: '2023-05-22T08:55:00Z' },
          total_records: { type: :integer, example: 100 },
          processed_records: { type: :integer, example: 45 },
          success_count: { type: :integer, example: 40 },
          failure_count: { type: :integer, example: 5 }
        }
      }
    end
  end
end
