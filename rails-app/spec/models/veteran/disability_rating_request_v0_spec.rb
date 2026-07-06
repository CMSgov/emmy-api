require 'rails_helper'

module Veteran
  RSpec.describe DisabilityRatingRequestV0, type: :model do
    describe '#to_va_payload' do
      it 'maps fields correctly' do
        params = {
          firstName: 'John',
          middleName: 'M',
          lastName: 'Doe',
          dateOfBirth: '1990-01-01',
          ssn: '999999999',
          address: {
            street1: '123 Main St',
            city: 'Anytown',
            state: 'NY',
            postalCode: '12345',
            country: 'USA'
          }
        }
        req = DisabilityRatingRequestV0.new(params)
        payload = req.to_va_payload

        expect(payload[:first_name]).to eq('John')
        expect(payload[:middle_name]).to eq('M')
        expect(payload[:last_name]).to eq('Doe')
        expect(payload[:birth_date]).to eq('1990-01-01')
        expect(payload[:ssn]).to eq('999999999')
        expect(payload[:street_address_line1]).to eq('123 Main St')
        expect(payload[:city]).to eq('Anytown')
        expect(payload[:state]).to eq('NY')
        expect(payload[:zipcode]).to eq('12345')
        expect(payload[:country]).to eq('USA')
      end

      it 'omits optional fields' do
        params = {
          firstName: 'John',
          lastName: 'Doe',
          dateOfBirth: '1990-01-01'
        }
        req = DisabilityRatingRequestV0.new(params)
        payload = req.to_va_payload

        expect(payload[:first_name]).to eq('John')
        expect(payload[:last_name]).to eq('Doe')
        expect(payload[:birth_date]).to eq('1990-01-01')
        expect(payload[:middle_name]).to be_nil
        expect(payload[:ssn]).to be_nil
        expect(payload[:street_address_line1]).to be_nil
      end
    end

    describe '#can_use_restricted_endpoint?' do
      it 'returns true when ssn is present' do
        req = DisabilityRatingRequestV0.new(ssn: '999999999')
        expect(req.can_use_restricted_endpoint?).to be true
      end

      it 'returns false when ssn is absent' do
        req = DisabilityRatingRequestV0.new
        expect(req.can_use_restricted_endpoint?).to be false
      end

      it 'returns false when ssn is empty' do
        req = DisabilityRatingRequestV0.new(ssn: '')
        expect(req.can_use_restricted_endpoint?).to be false
      end
    end
  end
end
