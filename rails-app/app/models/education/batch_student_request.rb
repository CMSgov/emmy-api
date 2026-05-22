module Education
  class BatchStudentRequest
    attr_accessor :batch_id, :submitted_by, :callback_url, :students

    def initialize(params = {})
      @batch_id = params[:batchId]
      @submitted_by = params[:submittedBy]
      @callback_url = params[:callbackUrl]
      @students = (params[:students] || []).map do |s|
        {
          record_id: s[:recordId],
          first_name: s[:firstName],
          last_name: s[:lastName],
          date_of_birth: s[:dateOfBirth],
          ssn: s[:ssn]
        }
      end
    end

    def self.to_swagger_schema
      {
        type: :object,
        required: %i[batchId submittedBy students],
        properties: {
          batchId: {
            type: :string,
            description: 'A unique identifier for the batch. This must be unique to any previous batch requested by this client.',
            example: 'batch-2023-05-22-001'
          },
          submittedBy: { type: :string, example: 'admin@university.edu' },
          callbackUrl: { type: :string, example: 'https://university.edu/api/callbacks/enrollment' },
          students: {
            type: :array,
            items: {
              type: :object,
              required: %i[recordId firstName lastName dateOfBirth],
              properties: {
                recordId: {
                  type: :string,
                  description: 'A unique identifier for this record within the batch.',
                  example: 'STUDENT-1001'
                },
                firstName: { type: :string, example: 'Jane' },
                lastName: { type: :string, example: 'Smith' },
                dateOfBirth: { type: :string, example: '1995-03-15' },
                ssn: { type: :string, example: '000-00-1111' }
              }
            }
          }
        }
      }
    end
  end
end
