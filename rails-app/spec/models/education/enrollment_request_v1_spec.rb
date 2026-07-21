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

    context 'when asOfDate is provided' do
      let(:terms_accepted) { true }

      it 'includes asOfDate in the payload' do
        expect(payload[:asOfDate]).to eq("2026-01-01")
      end
    end
  end

  describe '#ssn_only?' do
    it 'returns true when only SSN is provided' do
      request = described_class.new(personSocialSecurityNumber: '123456789')
      expect(request.ssn_only?).to be true
    end

    it 'returns false when name is provided' do
      request = described_class.new(personSocialSecurityNumber: '123456789', personGivenName: 'John')
      expect(request.ssn_only?).to be false
    end

    it 'returns false when surname is provided' do
      request = described_class.new(personSocialSecurityNumber: '123456789', personSurName: 'Doe')
      expect(request.ssn_only?).to be false
    end

    it 'returns false when birth date is provided' do
      request = described_class.new(personSocialSecurityNumber: '123456789', personBirthDate: '1990-01-01')
      expect(request.ssn_only?).to be false
    end

    it 'returns false when SSN is missing' do
      request = described_class.new(personGivenName: 'John')
      expect(request.ssn_only?).to be false
    end
  end
end
