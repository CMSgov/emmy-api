module Education
  class EnrollmentResponseV1
    attr_accessor :enrollment_status, :enrollment_details, :data_source, :metadata, :student_info_provided, :transaction_details

    def initialize(params = {})
      @enrollment_details = (params[:enrollmentDetails] || []).map do |detail|
        d = detail.dup
        d[:officialSchoolName] = d.delete(:schoolName) if d.key?(:schoolName)
        d
      end
      @data_source = params[:dataSource]
      @metadata = params[:metadata]
      @student_info_provided = params[:studentInfoProvided]
      @transaction_details = params[:transactionDetails]
    end

    def self.map_transaction_details(transaction_details)
      {
        transactionId: transaction_details["transactionId"],
        orderId: transaction_details["orderId"],
        transactionStatusCode: transaction_details["transactionStatusCode"],
        transactionFee: transaction_details["transactionFee"],
        salesTax: transaction_details["salesTax"],
        transactionTotal: transaction_details["transactionTotal"],
        requestedByText: transaction_details["requestedByText"],
        requestedDateTimeText: transaction_details["requestedDateTimeText"],
        notifiedDateTimeText: transaction_details["notifiedDateTimeText"],
        nscHitIndicator: transaction_details["nscHitIndicator"]
      }
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          enrollmentDetails: {
            type: :array,
            items: {
              type: :object,
              properties: {
                officialSchoolName: { type: :string, description: "Official name of the educational institution.", example: "University of Excellence" },
                schoolCode: { type: :string, description: "Six digit School ID assigned by the Department of Education for the educational istitution.", example: "001171" },
                branchCode: { type: :string, description: "Branch code for the school.", example: "00" },
                enrollmentData: {
                  type: :array,
                  items: {
                    type: :object,
                    properties: {
                      termBeginDate: { type: :string, description: "Start date of the academic term (YYYY-MM-DD).", example: "2023-01-15" },
                      termEndDate: { type: :string, description: "End date of the academic term (YYYY-MM-DD).", example: "2023-05-20" },
                      enrollmentStatusCode: { type: :string, description: "Enrollment status code for this specific term.", example: "F" },
                      schoolCertifiedOnDate: { type: :string, description: "The date when the school certified the enrollment.", example: "2023-01-20" },
                      anticipatedGraduationDate: { type: :string, description: "The anticipated graduation date.", example: "2026-05-20" }
                    }
                  }
                }
              }
            }
          },
          studentInfoProvided: {
            type: :object,
            properties: {
              personGivenName: { type: :string },
              personSurName: { type: :string },
              previousNames: {
                type: :array,
                items: {
                  type: :object,
                  properties: {
                    personGivenName: { type: :string },
                    personMiddleName: { type: :string },
                    personSurName: { type: :string }
                  }
                }
              },
              personBirthDate: { type: :string }
            }
          },
          transactionDetails: {
            type: :object,
            properties: {
              transactionId: { type: :string },
              orderId: { type: :string },
              transactionStatusCode: { type: :string },
              transactionFee: { type: :string },
              salesTax: { type: :string },
              transactionTotal: { type: :string },
              requestedByText: { type: :string },
              requestedDateTimeText: { type: :string },
              notifiedDateTimeText: { type: :string },
              nscHitIndicator: { type: :boolean }
            }
          },
          dataSource: { type: :string, description: "The source of the enrollment data (e.g., NSC).", example: "NSC" },
          responseMetadata: {
            type: :object,
            properties: {
              responseCode: { type: :string, example: "MS000000" },
              responseText: { type: :string, example: "Success" }
            }
          }
        }
      }
    end

    def as_json(options = {})
      {
        enrollmentDetails: enrollment_details,
        studentInfoProvided: student_info_provided,
        transactionDetails: self.class.map_transaction_details(transaction_details),
        dataSource: data_source,
        responseMetadata: {
          responseCode: "MS000000",
          responseText: "Success"
        }
      }
    end
  end
end
