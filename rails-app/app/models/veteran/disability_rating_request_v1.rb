module Veteran
  class DisabilityRatingRequestV1
    attr_accessor :person_social_security_number, :person_given_name, :person_middle_name, :person_sur_name,
                  :person_birth_date, :person_sex_code, :address

    def initialize(params = {})
      vadr = params[:vadrRequest] || {}
      @person_social_security_number = vadr[:personSocialSecurityNumber]
      @person_given_name = vadr[:personGivenName]
      @person_middle_name = vadr[:personMiddleName]
      @person_sur_name = vadr[:personSurName]
      @person_birth_date = vadr[:personBirthDate]
      @person_sex_code = vadr[:personSexCode]
      @address = vadr[:personContactInformation] || {}
    end

    def can_use_restricted_endpoint?
      person_social_security_number.present?
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          vadrRequest: {
            type: :object,
            required: %i[personGivenName personSurName personBirthDate],
            properties: {
              personGivenName: { type: :string, description: "First name of the veteran.", example: "John" },
              personMiddleName: { type: :string, description: "Middle name of the veteran.", example: "Quincy" },
              personSurName: { type: :string, description: "Last name of the veteran.", example: "Doe" },
              personBirthDate: { type: :string, description: "Date of birth of the veteran (YYYY-MM-DD).", example: "1990-01-01" },
              personSocialSecurityNumber: { type: :string, description: "Social Security Number of the veteran.", example: "000-00-0000" },
              personSexCode: { type: :string, description: "Gender of the veteran.", enum: %w[M F], example: "M" },
              personContactInformation: {
                type: :object,
                properties: {
                  streetLineOneAddress: { type: :string, example: "123 Main St" },
                  streetLineTwoAddress: { type: :string, example: "Apt 4B" },
                  cityName: { type: :string, example: "Arlington" },
                  stateText: { type: :string, example: "VA" },
                  zipCode: { type: :string, example: "22202" },
                  countryText: { type: :string, example: "USA" },
                  telephoneNumber: { type: :string, example: "555-555-5555" }
                }
              }
            }
          }
        }
      }
    end

    def to_va_payload
      out = {
        first_name: person_given_name,
        middle_name: person_middle_name,
        last_name: person_sur_name,
        birth_date: person_birth_date,
        ssn: person_social_security_number,
        gender: person_sex_code
      }.compact

      if address.present?
        out.merge!({
          street_address_line1: address[:streetLineOneAddress],
          street_address_line2: address[:streetLineTwoAddress],
          city: address[:cityName],
          state: address[:stateText],
          zipcode: address[:zipCode],
          country: address[:countryText],
          phone: address[:telephoneNumber]
        }.compact)
      end
      out
    end
  end
end
