class ApplicationController < ActionController::API
  before_action :set_current_request_id

  private

  def client_id
    sub = request.headers["X-Sub"]
    return sub if sub.present?

    # Try extracting from Bearer token in Authorization header
    auth_header = request.headers["Authorization"]
    if auth_header&.start_with?("Bearer ")
      token = auth_header.split(" ").last
      begin
        # We parse without verification to extract the 'sub' claim,
        # mirroring the Go SubjectMiddleware logic.
        decoded_token = JWT.decode(token, nil, false)
        return decoded_token.first["sub"] || "unknown-subject"
      rescue JWT::DecodeError
        return "unknown-subject"
      end
    end

    "unknown-subject"
  end

  def set_current_request_id
    Current.request_id = request.request_id
  end
end
