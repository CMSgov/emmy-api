module Veteran
  class DisabilityRatingRequest
    attr_accessor :first_name, :middle_name, :last_name, :date_of_birth, :ssn, :address

    def initialize(params = {})
      @first_name = params[:first_name]
      @middle_name = params[:middle_name]
      @last_name = params[:last_name]
      @date_of_birth = params[:date_of_birth]
      @ssn = params[:ssn]
      @address = params[:address]
    end

    def to_va_payload
      out = {
        first_name: first_name,
        middle_name: middle_name,
        last_name: last_name,
        birth_date: date_of_birth,
        ssn: ssn
      }.compact

      if address
        out.merge!({
          street_address_line1: address[:street1],
          street_address_line2: address[:street2],
          street_address_line3: address[:street3],
          city: address[:city],
          state: address[:state],
          zipcode: address[:postalCode],
          country: address[:country]
        }.compact)
      end

      out
    end
  end
end
