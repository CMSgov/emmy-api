require 'rails_helper'

RSpec.describe "Health", type: :request do
  describe "GET /health" do
    it "should get health" do
      redis_mock = instance_double(Redis)
      allow(redis_mock).to receive(:ping).and_return("PONG")
      allow(redis_mock).to receive(:close)
      allow(Redis).to receive(:new).and_return(redis_mock)

      get "/health"
      expect(response).to have_http_status(:success)
    end

    it "should return service_unavailable when redis is down" do
      redis_mock = instance_double(Redis)
      allow(redis_mock).to receive(:ping).and_raise(Redis::CannotConnectError, "Error connecting to Redis")
      allow(redis_mock).to receive(:close)
      allow(Redis).to receive(:new).and_return(redis_mock)

      get "/health"
      expect(response).to have_http_status(:service_unavailable)
    end
  end
end
