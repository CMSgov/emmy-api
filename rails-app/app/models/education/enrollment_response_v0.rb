module Education
  class EnrollmentResponseV0
    attr_accessor :enrollment_status, :enrollment_details, :raw_data, :data_source, :metadata

    def initialize(params = {})
      @enrollment_status = params[:enrollmentStatus]
      @enrollment_details = params[:enrollmentDetails] || []
      @raw_data = params[:rawData]
      @data_source = params[:dataSource]
      @metadata = params[:metadata]
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          enrollmentStatus: {
            type: :string,
            description: "Aggregated enrollment status across all schools.",
            enum: EnrollmentStatus::RANKS.keys,
            example: "FULL_TIME"
          },
          enrollmentDetails: {
            type: :array,
            items: {
              type: :object,
              properties: {
                schoolName: { type: :string, description: "Official name of the educational institution.", example: "University of Excellence" },
                schoolCode: { type: :string, description: "Six digit School ID assigned by the Department of Education for the educational istitution.", example: "001171" },
                termBeginDate: { type: :string, description: "Start date of the academic term (YYYY-MM-DD).", example: "2023-01-15" },
                termEndDate: { type: :string, description: "End date of the academic term (YYYY-MM-DD).", example: "2023-05-20" },
                enrollmentStatus: {
                  type: :string,
                  description: "Enrollment status for this specific term.",
                  enum: EnrollmentStatus::RANKS.keys,
                  example: "FULL_TIME"
                }
              }
            }
          },
          dataSource: { type: :string, description: "The source of the enrollment data (e.g., NSC).", example: "NSC" },
          metadata: {
            type: :object,
            properties: {
              durationMs: { type: :integer, description: "Time taken to fetch data from the source in milliseconds.", example: 125 }
            }
          }
        }
      }
    end

    def as_json(options = {})
      {
        enrollmentStatus: enrollment_status,
        enrollmentDetails: enrollment_details,
        rawData: raw_data,
        dataSource: data_source,
        metadata: metadata
      }
    end
  end
end
