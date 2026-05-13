class ApplicationController < ActionController::API
  before_action :set_current_request_id

  private

  def set_current_request_id
    Current.request_id = request.request_id
  end
end
