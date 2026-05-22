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
          batchId: { type: :string, example: 'test-batch-123' },
          submittedBy: { type: :string, example: 'user@example.com' },
          callbackUrl: { type: :string, example: 'https://example.com/callback' },
          students: {
            type: :array,
            items: {
              type: :object,
              required: %i[recordId firstName lastName dateOfBirth],
              properties: {
                recordId: { type: :string, example: 'rec1' },
                firstName: { type: :string, example: 'John' },
                lastName: { type: :string, example: 'Doe' },
                dateOfBirth: { type: :string, example: '1990-01-01' },
                ssn: { type: :string, example: '000-00-0000' }
              }
            }
          }
        }
      }
    end
  end
end
