require 'net/http'
require 'uri'
require 'json'
require 'base64'

# Configuration from environment variables
CLIENT_ID = ENV['CLIENT_ID']
SECRET_ID = ENV['SECRET_ID']
EMMY_ENV = ENV['EMMY_ENV'] || 'uat'

# Default mappings for known environments
ENV_DEFAULTS = {
  'dev' => {
    api_url: 'https://api.dev.emmy.cms.gov/api/v0',
    cognito_url: 'https://api.dev.emmy.cms.gov/oauth2/token'
  },
  'uat' => {
    api_url: 'https://api.uat.emmy.cms.gov/api/v0',
    cognito_url: 'https://emmy-uat.auth.us-east-1.amazoncognito.com/oauth2/token'
  },
  'demo' => {
    api_url: 'https://api.demo.emmy.cms.gov/api/v0',
    cognito_url: 'https://api.demo.emmy.cms.gov/oauth2/token'
  }
}.freeze

default_config = ENV_DEFAULTS[EMMY_ENV] || ENV_DEFAULTS['uat']

EMMY_API_COGNITO_URL = ENV['EMMY_API_COGNITO_URL'] || default_config[:cognito_url]
EMMY_API_URL = ENV['EMMY_API_URL'] || default_config[:api_url]

if CLIENT_ID.nil? || SECRET_ID.nil?
  puts "Error: CLIENT_ID or SECRET_ID environment variables not set."
  exit 1
end

def get_jwt(client_id, secret_id)
  uri = URI.parse(EMMY_API_COGNITO_URL)
  request = Net::HTTP::Post.new(uri)
  request.basic_auth(client_id, secret_id)
  request.set_form_data({ "grant_type" => "client_credentials" })
  request.content_type = "application/x-www-form-urlencoded"

  use_ssl = uri.scheme == 'https'
  options = { use_ssl: use_ssl }
  options[:verify_mode] = 0 if use_ssl

  response = Net::HTTP.start(uri.hostname, uri.port, options) do |http|
    http.request(request)
  end

  if response.code == "200"
    JSON.parse(response.body)['access_token']
  else
    puts "Failed to get JWT: #{response.code} #{response.body}"
    exit 1
  end
end

def post_request(url, jwt, body)
  uri = URI.parse(url)
  request = Net::HTTP::Post.new(uri)
  request["Authorization"] = "Bearer #{jwt}"
  request["Content-Type"] = "application/json"
  request.body = JSON.dump(body)

  # Note: -k in curl is equivalent to verify_mode: OpenSSL::SSL::VERIFY_NONE
  # Use with caution in production.
  use_ssl = uri.scheme == 'https'
  options = { use_ssl: use_ssl }
  options[:verify_mode] = 0 if use_ssl

  response = Net::HTTP.start(uri.hostname, uri.port, options) do |http|
    http.request(request)
  end

  puts "POST #{url}"
  puts "Response: #{response.code} #{response.body}"

  if response.code.to_i >= 400
    puts "Request failed with status #{response.code}"
    exit 1
  end
end

puts "Obtaining JWT..."
jwt = get_jwt(CLIENT_ID, SECRET_ID)
puts "JWT obtained successfully."

puts "\nTesting Education Enrollments..."
post_request(
  "#{EMMY_API_URL}/education-enrollments",
  jwt,
  {
    "firstName" => "Lynette",
    "lastName" => "Oyola",
    "dateOfBirth" => "1988-10-24"
  }
)

puts "\nTesting Veteran Disability Ratings..."
post_request(
  "#{EMMY_API_URL}/veteran-disability-ratings",
  jwt,
  {
    "firstName" => "Alfredo",
    "lastName" => "Armstrong",
    "dateOfBirth" => "1993-06-08",
    "ssn" => "796-01-2476"
  }
)

puts "\nAll tests completed."
