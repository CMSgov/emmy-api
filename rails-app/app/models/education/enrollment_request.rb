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

    def to_nsc_payload(account_id)
      out = {
        accountId: account_id,
        dateOfBirth: date_of_birth,
        lastName: last_name,
        firstName: first_name,
        middleName: middle_name,
        ssn: ssn,
        endClient: "CMS",
        terms: "y"
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
