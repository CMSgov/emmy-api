# Be sure to restart your server when you modify this file.

# Configure parameters to be partially matched (e.g. passw matches password) and filtered from the log file.
# Use this to limit dissemination of sensitive information.
# See the ActiveSupport::ParameterFilter documentation for supported notations and behaviors.
Rails.application.config.filter_parameters += [
  :passw, :email, :secret, :token, :_key, :crypt, :salt, :certificate, :otp, :ssn, :cvv, :cvc,
  :firstName, :lastName, :middleName, :first_name, :last_name, :middle_name,
  :dateOfBirth, :address, :phone, :ssn, :dob, :date_of_birth, :dateOfBirth, :dob, :birthdate, :birth_date, :birthDate
]
