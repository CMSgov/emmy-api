require 'swagger_helper'

RSpec.describe 'api/v0/education_enrollments', type: :request do

  path '/api/v0/education-enrollments' do

    before do
      fake_result = {
        enrollmentStatus: "FULL_TIME",
        dataSource: "NSC",
        metadata: {
          durationMs: 12
        }
      }
      coordinator_mock = instance_double(Education::ServiceCoordinator)
      expect(coordinator_mock).to receive(:lookup_enrollment_status)
                                    .with(
                                      {
                                        firstName: "John",
                                        lastName: "Doe",
                                        dateOfBirth: "1990-01-01"
                                      }
                                    )
                                    .and_return(fake_result)

      allow(Education::ServiceCoordinator).to receive(:new).and_return(coordinator_mock)
    end

    post('create education-enrollment') do
      security [oauth2: []]
      let(:Authorization) { 'Bearer <token>' }
      consumes 'application/json'
      produces 'application/json'
      parameter name: :education_enrollment,
                in: :body,
                schema: Education::EnrollmentRequest.to_swagger_schema

      response(200, 'successful') do
        schema Education::EnrollmentResponse.to_swagger_schema

        let(:education_enrollment) do
          {
            firstName: 'John',
            lastName: 'Doe',
            dateOfBirth: '1990-01-01'
          }
        end

        example 'application/json', :successful_enrollment_lookup, {
          enrollmentStatus: "FULL_TIME",
          enrollmentDetails: [
            {
              schoolName: "University of Excellence",
              termBeginDate: "2023-01-15",
              termEndDate: "2023-05-20",
              enrollmentStatus: "FULL_TIME"
            }
          ],
          dataSource: "NSC",
          metadata: {
            durationMs: 125
          }
        }

        after do |example|
          example.metadata[:response][:content] = {
            'application/json' => {
              example: JSON.parse(response.body, symbolize_names: true)
            }
          }
        end
        run_test!
      end
    end
  end
end
