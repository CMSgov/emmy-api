module Education
  class BatchDetailsRequest
    attr_accessor :batch_job_id

    def initialize(params = {})
      @batch_job_id = params[:batchJobId]
    end
  end
end
