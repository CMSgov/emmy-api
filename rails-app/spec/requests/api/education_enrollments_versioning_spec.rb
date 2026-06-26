require 'rails_helper'

RSpec.describe 'Education Enrollment Response Versioning', type: :request do
  let(:nsc_response) do
    {
      "status" => { "code" => "0", "message" => "Successful" },
      "studentInfoProvided" => {
        "dateOfBirth" => "1988-10-24",
        "firstName" => "Lynette",
        "lastName" => "Oyola",
        "middleName" => ""
      },
      "transactionDetails" => {
        "nscHit" => "Y",
        "transactionId" => "4001686568",
        "orderId" => "2002833179",
        "transactionStatus" => "CNF",
        "transactionFee" => "0.00",
        "salesTax" => "0.00",
        "transactionTotal" => "0.00",
        "requestedByText" => "Test Billing",
        "requestedDateTimeText" => "2026-04-23 14:20:41.049",
        "notifiedDateTimeText" => "2026-04-23 14:20:41.049"
      },
      "enrollmentDetails" => [
            {
              "officialSchoolName" => "Test University",
              "schoolCode" => "123456",
              "currentEnrollmentStatus" => "FT",
              "enrollmentData" => [
                {
                  "enrollmentStatus" => "F",
                  "termBeginDate" => "2023-01-01",
                  "termEndDate" => "2023-05-01"
                }
              ]
            }
          ]
        }
      end

      before do
        # Mock the NSC response at the client level
        allow_any_instance_of(Education::NscClient).to receive(:lookup_enrollment_status) do |_instance, enrollment_req|
          version = enrollment_req.is_a?(Education::EnrollmentRequestV1) ? :v1 : :v0
          Education::EnrollmentMapper.translate_nsc_response(enrollment_req, nsc_response, 125, version: version)
        end
      end

      describe 'V0 Response' do
        it 'returns schoolName in enrollmentDetails' do
          payload = {
            firstName: 'John',
            lastName: 'Doe',
            dateOfBirth: '1990-01-01'
          }

          post '/api/v0/education-enrollments', params: payload, as: :json

          expect(response).to have_http_status(:ok)
          json = JSON.parse(response.body, symbolize_names: true)

          expect(json[:enrollmentDetails].first).to have_key(:schoolName)
          expect(json[:enrollmentDetails].first[:schoolName]).to eq("Test University")
          expect(json[:enrollmentDetails].first).not_to have_key(:officialSchoolName)
          expect(json).to have_key(:rawData)
        end
      end

      describe 'V1 Response' do
        it 'returns officialSchoolName in enrollmentDetails' do
          payload = {
            personSocialSecurityNumber: "123456789",
            personGivenName: "John",
            personSurName: "Doe",
            personBirthDate: "1988-02-29"
          }

          post '/api/v1/education-enrollments', params: payload, as: :json

          expect(response).to have_http_status(:ok)
          json = JSON.parse(response.body, symbolize_names: true)

          expect(json).to have_key(:nscResponse)
          nsc_resp = json[:nscResponse]

          expect(nsc_resp[:enrollmentDetails].first).to have_key(:officialSchoolName)
          expect(nsc_resp[:enrollmentDetails].first[:officialSchoolName]).to eq("Test University")
          expect(nsc_resp[:enrollmentDetails].first).to have_key(:currentEnrollmentStatusCode)
          expect(nsc_resp[:enrollmentDetails].first[:currentEnrollmentStatusCode]).to eq("FT")
          expect(nsc_resp[:enrollmentDetails].first).not_to have_key(:schoolName)
          expect(nsc_resp).not_to have_key(:rawData)
      expect(nsc_resp).to have_key(:studentInfoProvided)
      expect(nsc_resp[:studentInfoProvided][:personGivenName]).to eq("John")
      expect(nsc_resp[:studentInfoProvided][:personSurName]).to eq("Doe")
      expect(nsc_resp[:studentInfoProvided][:personBirthDate]).to eq("1988-02-29")
      expect(nsc_resp).to have_key(:transactionDetails)
      expect(nsc_resp[:transactionDetails][:transactionId]).to eq("4001686568")
      expect(nsc_resp[:transactionDetails][:transactionStatusCode]).to eq("CNF")
      expect(nsc_resp[:transactionDetails][:nscHitIndicator]).to eq("Y")
    end
  end
end
