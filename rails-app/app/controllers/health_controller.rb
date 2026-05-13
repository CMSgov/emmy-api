class HealthController < ApplicationController
  REDIS_ADDR = ENV.fetch("REDIS_ADDR") { "localhost:6379" }
  REDIS_DB = ENV.fetch("REDIS_DB") { "1" }
  REDIS_PASSWORD = ENV.fetch("REDIS_PASSWORD", nil)
  REDIS_SSL_URL = ENV.fetch("REDIS_SSL_URL", "redis")

  if REDIS_PASSWORD.present?
    REDIS_URL = "#{REDIS_SSL_URL}://:#{REDIS_PASSWORD}@#{REDIS_ADDR}/#{REDIS_DB}"
  else
    REDIS_URL = "#{REDIS_SSL_URL}://#{REDIS_ADDR}/#{REDIS_DB}"
  end

  if REDIS_SSL_URL == "rediss"
    VERIFY_MODE = OpenSSL::SSL::VERIFY_PEER
  else
    VERIFY_MODE = OpenSSL::SSL::VERIFY_NONE
  end

  def show
    redis = Redis.new(url: REDIS_URL,  ssl_params: {
      verify_mode: VERIFY_MODE
    })

    begin
      redis.ping
      head :ok
    rescue => e
      Rails.logger.error "Health check failed: #{e.message} #{REDIS_ADDR}"
      head :service_unavailable
    ensure
      redis.close if redis
    end
  end
end
