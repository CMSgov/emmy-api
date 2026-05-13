module Veteran
  class DisabilityRatingMapper
    def self.map_response(total_disability_response, datasource_duration, start_time)
      response_timestamp = Time.now

      total_disabilty_data = total_disability_response.dig("data")
      Veteran::DisabilityRatingResponse.new(
        totalDisabilityStatus: total_disabilty_data["total_disability"]["status"],
        totalDisabilityStatusEffectiveDate: total_disabilty_data["total_disability"]["effective_date"],
        rawData: total_disability_response,
        dataSource: "VA",
        metadata: {
          apiVersion: ENV['SERVICE_VERSION'] || '1.3.0',
          environment: ENV['ENVIRONMENT'] || 'development',
          requestTimestamp: start_time.utc.iso8601(3),
          responseTimestamp: response_timestamp.utc.iso8601(3),
          datasourceDurationMillis: datasource_duration.to_i,
          transactionId: Current.request_id || SecureRandom.uuid
        }
      )
    end
  end
end
