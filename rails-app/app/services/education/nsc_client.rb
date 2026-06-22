require "net/http"

module Education
  class NotFoundError < StandardError; end

  class NscClient
    def initialize
    end

    def lookup_enrollment_status(enrollment_req)
      token = fetch_oauth_token
      nsc_payload = enrollment_req.to_nsc_payload(ENV["NSC_ACCOUNT_ID"])

      uri = URI(ENV["NSC_SUBMIT_URL"])
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = (uri.scheme == "https")

      request = Net::HTTP::Post.new(uri.path, {
        "Content-Type" => "application/json",
        "Accept" => "application/json",
        "Authorization" => "Bearer #{token}"
      })
      request.body = nsc_payload.to_json

      start_time = Time.now
      response = http.request(request)
      duration = (Time.now - start_time) * 1000


      if response.code.to_i < 200 || response.code.to_i >= 300
        raise "NSC submit failed: status=#{response.code} body=#{response.body} correlationId=#{enrollment_req.correlation_id}"
      end

      nsc_resp = JSON.parse(response.body)
      Education::EnrollmentMapper.translate_nsc_response(nsc_resp, duration)
    end

    private

    def fetch_oauth_token
      uri = URI(ENV["NSC_TOKEN_URL"])
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = (uri.scheme == "https")

      request = Net::HTTP::Post.new(uri.path)
      request.basic_auth(ENV["NSC_CLIENT_ID"], ENV["NSC_CLIENT_SECRET"])
      request.set_form_data({
        grant_type: "client_credentials",
        scope: "vs.api.insights"
      })

      response = http.request(request)

      if response.code.to_i < 200 || response.code.to_i >= 300
        raise "NSC OAuth failed: status=#{response.code} body=#{response.body}"
      end

      JSON.parse(response.body)["access_token"]
    end
  end
end
