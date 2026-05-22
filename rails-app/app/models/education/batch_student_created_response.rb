module Education
  class BatchStudentCreatedResponse
    attr_accessor :message, :batch_job_id

    def initialize(message:, batch_job_id:)
      @message = message
      @batch_job_id = batch_job_id
    end

    def as_json(options = {})
      {
        message: message,
        batchJobId: batch_job_id
      }
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          message: { type: :string, example: 'Batch registration successful' },
          batchJobId: { type: :string, example: 'test-batch-123' }
        }
      }
    end
  end
end
