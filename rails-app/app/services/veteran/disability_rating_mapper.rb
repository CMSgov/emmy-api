module Veteran
  class DisabilityRatingMapper
    def self.map_response(va_resp, datasource_duration, start_time)
      attributes = va_resp.dig('data', 'attributes') || {}

      individual_ratings = attributes['individual_ratings'] || []
      earliest_end_date = individual_ratings.map { |r| r['rating_end_date'] }
                                            .reject { |d| d.blank? }
                                            .min

      response_timestamp = Time.now

      Veteran::DisabilityRatingResponse.new(
        combinedDisabilityRating: attributes['combined_disability_rating'],
        combinedEffectiveDate: attributes['combined_effective_date'],
        legalEffectiveDate: attributes['legal_effective_date'],
        earliestRatingEndDate: earliest_end_date,
        rawData: va_resp,
        dataSource: "VA",
        metadata: {
          apiVersion: ENV['SERVICE_VERSION'] || '1.3.0',
          environment: ENV['ENVIRONMENT'] || 'development',
          requestTimestamp: start_time.utc.iso8601(3),
          responseTimestamp: response_timestamp.utc.iso8601(3),
          datasourceDurationMillis: datasource_duration.to_i,
          transactionId: SecureRandom.uuid
        }
      )
    end
  end
end
