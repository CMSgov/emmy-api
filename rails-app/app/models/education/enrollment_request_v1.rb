module Education
  class EnrollmentRequestV1
    attr_accessor :person_social_security_number, :person_given_name, :person_middle_name, :person_sur_name, :person_birth_date, :as_of_date, :previous_names, :terms_accepted_indicator

    def initialize(params = {})
      @person_social_security_number = params[:personSocialSecurityNumber]
      @person_given_name = params[:personGivenName]
      @person_middle_name = params[:personMiddleName]
      @person_sur_name = params[:personSurName]
      @person_birth_date = params[:personBirthDate]
      @as_of_date = params[:asOfDate]
      @previous_names = params[:previousNames] || []
      @terms_accepted_indicator = params[:termsAcceptedIndicator]
    end

    def correlation_id
      @correlation_id ||= SecureRandom.uuid
    end

    def to_nsc_payload(account_id)
      out = {
        accountId: account_id,
        dateOfBirth: person_birth_date,
        lastName: person_sur_name,
        firstName: person_given_name,
        middleName: person_middle_name,
        ssn: person_social_security_number,
        endClient: "CMS",
        terms: terms_accepted_indicator ? "y" : "n",
        correlationId: correlation_id,
        asOfDate: as_of_date
      }

      unless previous_names.empty?
        out[:previousNames] = previous_names.map do |pn|
          {
            firstName: pn[:personGivenName],
            middleName: pn[:personMiddleName],
            lastName: pn[:personSurName]
          }
        end
      end

      out
    end
  end
end
