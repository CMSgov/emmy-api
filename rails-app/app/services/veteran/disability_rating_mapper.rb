module Veteran
  class DisabilityRatingMapper
    def self.map_response(total_disability_response, disability_score_response, datasource_duration, start_time, version = "V0")
      response_timestamp = Time.now
      combined_data = disability_score_response.dig("data", "attributes") || {}

      individual_ratings = combined_data["individual_ratings"] || []
      earliest_end_date = individual_ratings.map { |r| r["rating_end_date"] }
                                            .reject { |d| d.blank? }
                                            .min

      total_disability_data = total_disability_response.dig("data") || {}
      if version == "V0"
        Veteran::DisabilityRatingResponseV0.new(
          permanentDisabilityStatus: total_disability_data.dig("permanent_and_total", "service_connected_status"),
          totalDisabilityStatus: total_disability_data.dig("total_disability", "status"),
          totalDisabilityStatusEffectiveDate: total_disability_data.dig("total_disability", "effective_date"),
          combinedDisabilityRating: combined_data["combined_disability_rating"],
          combinedEffectiveDate: combined_data["combined_effective_date"],
          legalEffectiveDate: combined_data["legal_effective_date"],
          rawData: { total_disability_response: total_disability_response, disability_score_response: disability_score_response },
          earliestRatingEndDate: earliest_end_date,
          dataSource: "VA",
          metadata: {
            apiVersion: ENV["SERVICE_VERSION"] || "1.3.0",
            environment: ENV["ENVIRONMENT"] || "development",
            requestTimestamp: start_time.utc.iso8601(3),
            responseTimestamp: response_timestamp.utc.iso8601(3),
            datasourceDurationMillis: datasource_duration.to_i,
            transactionId: Current.request_id || SecureRandom.uuid
          }
        )
      else
        Veteran::DisabilityRatingResponseV1.new(
          permanentDisabilityStatus: total_disability_data.dig("permanent_and_total", "service_connected_status"),
          totalDisabilityStatus: total_disability_data.dig("total_disability", "status"),
          totalDisabilityStatusEffectiveDate: total_disability_data.dig("total_disability", "effective_date"),
          combinedDisabilityRating: combined_data["combined_disability_rating"],
          combinedEffectiveDate: combined_data["combined_effective_date"],
          legalEffectiveDate: combined_data["legal_effective_date"],
          individualRatings: individual_ratings.map { |r|
            {
              decisionText: r["decision"],
              ratingEffectiveDate: r["effective_date"],
              ratingEndDate: r["rating_end_date"],
              ratingPercentage: r["rating_percentage"],
              disabilityRatingId: r["disability_rating_id"],
              staticIndicator: r["static_ind"]
            }
          },
          pensionAwardStatusIndicator: total_disability_data.dig("permanent_and_total", "pension_award_status"),
          rawData: { total_disability_response: total_disability_response, disability_score_response: disability_score_response },
          earliestRatingEndDate: earliest_end_date,
          dataSource: "VA",
          metadata: {
            apiVersion: ENV["SERVICE_VERSION"] || "1.3.0",
            environment: ENV["ENVIRONMENT"] || "development",
            requestTimestamp: start_time.utc.iso8601(3),
            responseTimestamp: response_timestamp.utc.iso8601(3),
            datasourceDurationMillis: datasource_duration.to_i,
            transactionId: Current.request_id || SecureRandom.uuid
          }
        )
      end
    end
  end
end
