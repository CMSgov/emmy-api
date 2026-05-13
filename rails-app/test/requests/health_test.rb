require "test_helper"

class HealthTest < ActionDispatch::IntegrationTest
  test "should get health" do
    # Mock Redis to ensure it doesn't try to connect to a real server during tests
    # if it's not available, although localhost:6379 might be available.
    # To be safe and deterministic:
    redis_mock = Minitest::Mock.new
    redis_mock.expect :ping, "PONG"
    redis_mock.expect :close, nil

    Redis.stub :new, redis_mock do
      get "/health"
      assert_response :success
    end

    redis_mock.verify
  end

  test "should return service_unavailable when redis is down" do
    redis_mock = Minitest::Mock.new
    redis_mock.expect :ping, nil do
      raise Redis::CannotConnectError, "Error connecting to Redis"
    end
    redis_mock.expect :close, nil

    Redis.stub :new, redis_mock do
      get "/health"
      assert_response :service_unavailable
    end

    redis_mock.verify
  end
end
