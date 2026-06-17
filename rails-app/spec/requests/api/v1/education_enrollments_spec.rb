require 'swagger_helper'

RSpec.describe 'api/v1/education_enrollments', type: :request do
  path '/api/v1/education-enrollments' do
    describe 'POST /api/v1/education-enrollments' do
      let(:fake_result) do
        {
          enrollmentStatus: "FULL_TIME",
          dataSource: "NSC",
          metadata: {
            durationMs: 12
          }
        }
      end

      before do
        coordinator_mock = instance_double(Education::ServiceCoordinator)
        allow(Education::ServiceCoordinator).to receive(:new).and_return(coordinator_mock)
        allow(coordinator_mock).to receive(:lookup_enrollment_status).with(kind_of(Education::EnrollmentRequestV1)).and_return(fake_result)
      end

      post('create education-enrollment') do
        security [ oauth2: [] ]
        let(:Authorization) { 'Bearer <token>' }
        consumes 'application/json'
        produces 'application/json'
        parameter name: :education_enrollment,
                  in: :body,
                  schema: {
                    type: :object,
                    properties: {
                      personSocialSecurityNumber: { type: :string },
                      personGivenName: { type: :string },
                      personMiddleName: { type: :string },
                      personSurName: { type: :string },
                      personBirthDate: { type: :string },
                      asOfDate: { type: :string },
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
                      }
                    }
                  }

        response(200, 'successful with new format') do
          let(:education_enrollment) do
            {
              personSocialSecurityNumber: "123456789",
              personGivenName: "John",
              personMiddleName: "Joe",
              personSurName: "Doe",
              previousNames: [
                {
                  personGivenName: "John",
                  personMiddleName: "Jacob",
                  personSurName: "Smith"
                },
                {
                  personGivenName: "Richard",
                  personSurName: "Roe"
                }
              ],
              personBirthDate: "1988-02-29",
              asOfDate: "2026-01-01"
            }
          end

          run_test!
        end
      end
    end
  end
end
