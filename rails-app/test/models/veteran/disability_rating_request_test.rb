require 'test_helper'

module Veteran
  class DisabilityRatingRequestTest < ActiveSupport::TestCase
    test 'to_va_payload maps fields correctly' do
      params = {
        first_name: 'John',
        middle_name: 'M',
        last_name: 'Doe',
        date_of_birth: '1990-01-01',
        ssn: '999999999',
        address: {
          street1: '123 Main St',
          city: 'Anytown',
          state: 'NY',
          postalCode: '12345',
          country: 'USA'
        }
      }
      req = DisabilityRatingRequest.new(params)
      payload = req.to_va_payload

      assert_equal 'John', payload[:first_name]
      assert_equal 'M', payload[:middle_name]
      assert_equal 'Doe', payload[:last_name]
      assert_equal '1990-01-01', payload[:birth_date]
      assert_equal '999999999', payload[:ssn]
      assert_equal '123 Main St', payload[:street_address_line1]
      assert_equal 'Anytown', payload[:city]
      assert_equal 'NY', payload[:state]
      assert_equal '12345', payload[:zipcode]
      assert_equal 'USA', payload[:country]
    end

    test 'to_va_payload omits optional fields' do
      params = {
        first_name: 'John',
        last_name: 'Doe',
        date_of_birth: '1990-01-01'
      }
      req = DisabilityRatingRequest.new(params)
      payload = req.to_va_payload

      assert_equal 'John', payload[:first_name]
      assert_equal 'Doe', payload[:last_name]
      assert_equal '1990-01-01', payload[:birth_date]
      assert_nil payload[:middle_name]
      assert_nil payload[:ssn]
      assert_nil payload[:street_address_line1]
    end
  end
end
