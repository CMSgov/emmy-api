module Veteran
  class DisabilityRatingRequest
    attr_accessor :first_name, :middle_name, :last_name, :date_of_birth, :ssn, :address

    def initialize(params = {})
      @first_name = params[:firstName]
      @middle_name = params[:middleName]
      @last_name = params[:lastName]
      @date_of_birth = params[:dateOfBirth]
      @ssn = params[:ssn]
      @address = params[:address]
    end

    def self.to_swagger_schema
      {
        type: :object,
        required: %i[firstName dateOfBirth lastName],
        properties: {
          firstName: { type: :string, description: "First name of the veteran.", example: "John" },
          middleName: { type: :string, description: "Middle name of the veteran.", example: "Quincy" },
          lastName: { type: :string, description: "Last name of the veteran.", example: "Doe" },
          dateOfBirth: { type: :string, description: "Date of birth of the veteran (YYYY-MM-DD).", example: "1990-01-01" },
          ssn: { type: :string, description: "Social Security Number of the veteran.", example: "000-00-0000" },
          address: {
            type: :object,
            description: "Postal address when a lookup requires demographic matching instead of SSN matching.",
            properties: {
              street1: { type: :string, description: "Primary street address line.", example: "123 Main St" },
              street2: { type: :string, description: "Secondary street address line when available.", example: "Apt 4B" },
              street3: { type: :string, description: "Additional street address line when available." },
              city: { type: :string, description: "City or locality.", example: "Arlington" },
              state: { type: :string, description: "State, province, or region code.", example: "VA" },
              postalCode: { type: :string, description: "Postal or ZIP code.", example: "22202" },
              country: { type: :string, description: "Country name or code.", example: "USA" }
            }
          }
        }
      }
    end

    def can_use_restricted_endpoint?
      ssn.present?
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
