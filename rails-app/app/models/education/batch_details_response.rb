module Education
  class BatchDetailsResponse
    attr_accessor :batch_job_id, :results

    def initialize(params = {})
      @batch_job_id = params[:batch_job_id]
      @results = params[:results] || []
    end

    def as_json(options = {})
      {
        batch_job_id: batch_job_id,
        results: results
      }
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          batch_job_id: { type: :string, example: "batch-2023-05-22-001" },
          results: {
            type: :array,
            items: {
              type: :object,
              properties: {
                record_id: { type: :string, example: "STUDENT-1001" },
                status: {
                  type: :string,
                  enum: Education::BatchStudent.statuses.values,
                  example: "SUCCESS"
                },
                found_enrollment: { type: :boolean, example: true },
                results: {
                  type: :array,
                  items: {
                    type: :object,
                    properties: {
                      schoolName: { type: :string, example: "University of Excellence" },
                      termBeginDate: { type: :string, example: "2023-01-15" },
                      termEndDate: { type: :string, example: "2023-05-20" },
                      enrollmentStatus: { type: :string, example: "FULL_TIME" }
                    }
                  }
                }
              }
            }
          }
        }
      }
    end
  end
end
