require 'rails_helper'

RSpec.describe Education::EnrollmentRequestV1, type: :model do
  describe '#to_nsc_payload' do
    let(:params) do
      {
        personSocialSecurityNumber: "123456789",
        personGivenName: "John",
        personSurName: "Doe",
        personBirthDate: "1988-02-29",
        asOfDate: "2026-01-01",
        termsAcceptedIndicator: terms_accepted
      }
    end
    let(:request) { described_class.new(params) }
    let(:account_id) { "test-account" }
    let(:payload) { request.to_nsc_payload(account_id) }

    context 'when termsAcceptedIndicator is true' do
      let(:terms_accepted) { true }

      it 'maps terms to "y"' do
        expect(payload[:terms]).to eq('y')
      end
    end

    context 'when termsAcceptedIndicator is false' do
      let(:terms_accepted) { false }

      it 'maps terms to "n"' do
        expect(payload[:terms]).to eq('n')
      end
    end

    context 'when termsAcceptedIndicator is nil' do
      let(:terms_accepted) { nil }

      it 'maps terms to "n"' do
        expect(payload[:terms]).to eq('n')
      end
    end
  end
end
