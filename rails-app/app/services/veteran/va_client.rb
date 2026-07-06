require "net/http"
require "openssl"
require "securerandom"

module Veteran
  class NotFoundError < StandardError; end

  class VaClient
    DISABILITY_RATING_PATH = "/disability_rating".freeze
    RESTRICTED_DISABILITY_RATING_PATH = "/restricted/disability_rating".freeze
    DISABILITY_RATING_SCOPE = "disability_rating_restricted.read disability_rating.read permanent_and_total_disability.read permanent_and_total_disability_restricted.read".freeze
    CLIENT_ASSERTION_TYPE = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer".freeze
    TOTAL_DISABILITY_PATH = "/permanent_and_total_disability"
    TOTAL_DISABILITY_RESTRICTED_PATH = "/restricted/permanent_and_total_disability"

    def initialize
    end


    def lookup_disability_rating(req_params, version = "V0")
      start_time = Time.now
      token = fetch_oauth_token

      # Ensure req_params is a request object
      rating_req = if req_params.respond_to?(:to_va_payload)
                     req_params
      else
                     Veteran::DisabilityRatingRequestV0.new(req_params)
      end

      total_disability_response = get_total_disability_response(rating_req, token)

      disability_score_response = get_total_disability_response(rating_req, token, restricted_endpoint: RESTRICTED_DISABILITY_RATING_PATH, unrestricted_endpoint: DISABILITY_RATING_PATH)
      datasource_duration = (Time.now - start_time) * 1000

      Veteran::DisabilityRatingMapper.map_response(JSON.parse(total_disability_response.body), JSON.parse(disability_score_response.body), datasource_duration, start_time, version)
    end

    private
    def get_total_disability_response(rating_req, token, restricted_endpoint: TOTAL_DISABILITY_RESTRICTED_PATH, unrestricted_endpoint: TOTAL_DISABILITY_PATH)
      va_payload = rating_req.to_va_payload
      if rating_req.can_use_restricted_endpoint?
        total_disability_uri = URI(File.join(ENV["VA_BASE_URL"], restricted_endpoint))
      else
        total_disability_uri = URI(File.join(ENV["VA_BASE_URL"], unrestricted_endpoint))
      end
      http = Net::HTTP.new(total_disability_uri.host, total_disability_uri.port)
      http.use_ssl = (total_disability_uri.scheme == "https")

      total_disability_request = Net::HTTP::Post.new(total_disability_uri.path, {
        "Content-Type" => "application/json",
        "Accept" => "application/json",
        "Authorization" => "Bearer #{token}"
      })

      total_disability_request.body = va_payload.to_json

      total_disability_response = http.request(total_disability_request)
      if total_disability_response.code.to_i == 404
        raise NotFoundError, "veteran not found"
      end

      if total_disability_response.code.to_i < 200 || total_disability_response.code.to_i >= 300
        Rails.logger.error("VA disability rating failed: status=#{total_disability_response.code} body=#{total_disability_response.body.to_s[0..800]}")
        raise "VA disability rating failed: status=#{total_disability_response.code}"
      end

      total_disability_response
    end

    def fetch_oauth_token
      assertion = signed_client_assertion

      uri = URI(ENV["VA_TOKEN_URL"])
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = (uri.scheme == "https")

      request = Net::HTTP::Post.new(uri.path)
      request.set_form_data({
        grant_type: "client_credentials",
        client_assertion_type: CLIENT_ASSERTION_TYPE,
        client_assertion: assertion,
        scope: DISABILITY_RATING_SCOPE
      })

      response = http.request(request)

      if response.code.to_i < 200 || response.code.to_i >= 300
        Rails.logger.error("VA OAuth failed: status=#{response.code} body=#{response.body}")
        raise "VA OAuth failed: status=#{response.code}"
      end

      JSON.parse(response.body)["access_token"]
    end

    def signed_client_assertion
      private_key = OpenSSL::PKey::RSA.new(File.read(ENV["VA_PRIVATE_KEY_PATH"]))

      now = Time.now.to_i
      payload = {
        iss: ENV["VA_CLIENT_ID"],
        sub: ENV["VA_CLIENT_ID"],
        aud: ENV["VA_AUD"],
        jti: SecureRandom.uuid,
        iat: now,
        exp: now + 60
      }

      JWT.encode(payload, private_key, "RS256")
    end
  end
end
