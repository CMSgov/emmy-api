require 'net/http'
require 'openssl'
require 'securerandom'

module Veteran
  class NotFoundError < StandardError; end

  class VaClient
    DISABILITY_RATING_PATH = "/restricted/disability_rating".freeze
    DISABILITY_RATING_SCOPE = "disability_rating_restricted.read".freeze
    CLIENT_ASSERTION_TYPE = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer".freeze

    def initialize
    end

    def lookup_disability_rating(req_params)
      start_time = Time.now

      rating_req = Veteran::DisabilityRatingRequest.new(req_params)
      token = fetch_oauth_token
      va_payload = rating_req.to_va_payload

      uri = URI(File.join(ENV['VA_BASE_URL'], DISABILITY_RATING_PATH))
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = (uri.scheme == 'https')

      request = Net::HTTP::Post.new(uri.path, {
        'Content-Type' => 'application/json',
        'Accept' => 'application/json',
        'Authorization' => "Bearer #{token}"
      })
      request.body = va_payload.to_json

      datasource_start_time = Time.now
      response = http.request(request)
      datasource_duration = (Time.now - datasource_start_time) * 1000

      if response.code.to_i == 404
        raise NotFoundError, "veteran not found"
      end

      if response.code.to_i < 200 || response.code.to_i >= 300
        Rails.logger.error("VA disability rating failed: status=#{response.code} body=#{response.body.to_s[0..800]}")
        raise "VA disability rating failed: status=#{response.code}"
      end

      va_resp = JSON.parse(response.body)
      Veteran::DisabilityRatingMapper.map_response(va_resp, datasource_duration, start_time)
    end

    private

    def fetch_oauth_token
      assertion = signed_client_assertion

      uri = URI(ENV['VA_TOKEN_URL'])
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = (uri.scheme == 'https')

      request = Net::HTTP::Post.new(uri.path)
      request.set_form_data({
        grant_type: 'client_credentials',
        client_assertion_type: CLIENT_ASSERTION_TYPE,
        client_assertion: assertion,
        scope: DISABILITY_RATING_SCOPE
      })

      response = http.request(request)

      if response.code.to_i < 200 || response.code.to_i >= 300
        Rails.logger.error("VA OAuth failed: status=#{response.code} body=#{response.body}")
        raise "VA OAuth failed: status=#{response.code}"
      end

      JSON.parse(response.body)['access_token']
    end

    def signed_client_assertion
      private_key = OpenSSL::PKey::RSA.new(File.read(ENV['VA_PRIVATE_KEY_PATH']))

      now = Time.now.to_i
      payload = {
        iss: ENV['VA_CLIENT_ID'],
        sub: ENV['VA_CLIENT_ID'],
        aud: ENV['VA_AUD'],
        jti: SecureRandom.uuid,
        iat: now,
        exp: now + 60
      }

      JWT.encode(payload, private_key, 'RS256')
    end

  end
end
