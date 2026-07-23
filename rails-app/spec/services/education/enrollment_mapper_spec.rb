require 'rails_helper'

RSpec.describe Education::EnrollmentMapper do
  before do
    # Ensure NotFoundError is defined
    unless defined?(Education::NotFoundError)
      module Education
        class NotFoundError < StandardError; end
      end
    end
  end
  let(:nsc_request) { double("nsc_request", person_given_name: "John", person_middle_name: "Q", person_sur_name: "Doe", previous_names: [], person_birth_date: "1988-02-29") }
  let(:duration) { 100 }

  describe '.translate_nsc_response' do
    context 'when nsc response is a no-hit' do
      let(:nsc_resp) do
        {
          "transactionDetails" => { "nscHit" => "N" },
          "enrollmentDetails" => []
        }
      end

      it 'raises NotFoundError for version v0' do
        expect {
          described_class.translate_nsc_response(nsc_request, nsc_resp, duration, version: :v0)
        }.to raise_error(Education::NotFoundError)
      end

      it 'does NOT raise NotFoundError for version v1 and returns response' do
        response = nil
        expect {
          response = described_class.translate_nsc_response(nsc_request, nsc_resp, duration, version: :v1)
        }.not_to raise_error
        expect(response.enrollment_details).to be_empty
      end
    end

    context 'when nsc response is not currently enrolled' do
      let(:nsc_resp) do
        {
          "transactionDetails" => { "nscHit" => "Y" },
          "enrollmentDetails" => [
            { "currentEnrollmentStatus" => "CN", "enrollmentData" => [] }
          ]
        }
      end

      it 'does NOT raise NotFoundError for version v0' do
        expect {
          described_class.translate_nsc_response(nsc_request, nsc_resp, duration, version: :v0)
        }.not_to raise_error
      end

      it 'does NOT raise NotFoundError for version v1 and returns response' do
        response = nil
        expect {
          response = described_class.translate_nsc_response(nsc_request, nsc_resp, duration, version: :v1)
        }.not_to raise_error
        expect(response.enrollment_details).not_to be_empty
      end
    end
  end
end
