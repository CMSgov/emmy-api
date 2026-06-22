module Education
  class EnrollmentRequest
    attr_accessor :first_name, :middle_name, :last_name, :date_of_birth, :ssn, :address

    def initialize(params = {})
      @first_name = params[:firstName]
      @middle_name = params[:middleName]
      @last_name = params[:lastName]
      @date_of_birth = params[:dateOfBirth]
      @ssn = params[:ssn]
      @address = params[:address]
    end

    def missing_required_field?
      first_name.blank? || last_name.blank? || date_of_birth.blank?
    end

    def self.to_swagger_schema
      {
        type: :object,
        required: %i[firstName dateOfBirth lastName],
        properties: {
          firstName: { type: :string, description: "First name of the student.", example: "John" },
          middleName: { type: :string, description: "Middle name of the student.", example: "Quincy" },
          lastName: { type: :string, description: "Last name of the student.", example: "Doe" },
          dateOfBirth: { type: :string, description: "Date of birth of the student (YYYY-MM-DD).", example: "1990-01-01" },
          ssn: { type: :string, description: "Social Security Number of the student.", example: "000-00-0000" },
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
            },
            required: %i[street1 city state postalCode country]
          }
        }
      }
    end

    def correlation_id
      @correlation_id ||= SecureRandom.uuid
    end

    def to_nsc_payload(account_id)
      out = {
        accountId: account_id,
        dateOfBirth: date_of_birth,
        lastName: last_name,
        firstName: first_name,
        middleName: middle_name,
        ssn: ssn,
        endClient: "CMS",
        terms: "y",
        correlationId: correlation_id
      }

      if address
        out[:address1] = address[:street1]
        out[:address2] = address[:street2]
        out[:city] = address[:city]
        out[:state] = address[:state]
        out[:zipCode] = address[:postalCode]
      end

      out
    end
  end
end
