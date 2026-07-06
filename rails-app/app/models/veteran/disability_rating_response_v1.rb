module Veteran
  class DisabilityRatingResponseV1
    attr_accessor :total_disability_status, :total_disability_status_effective_date,
                  :combined_disability_rating, :combined_effective_date, :legal_effective_date,
                  :earliest_rating_end_date, :permanent_disability_status, :metadata,
                  :pension_award_status_indicator, :individual_ratings

    def initialize(params = {})
      @total_disability_status = params[:totalDisabilityStatus]
      @total_disability_status_effective_date = params[:totalDisabilityStatusEffectiveDate]
      @combined_disability_rating = params[:combinedDisabilityRating]
      @combined_effective_date = params[:combinedEffectiveDate]
      @legal_effective_date = params[:legalEffectiveDate]
      @earliest_rating_end_date = params[:earliestRatingEndDate]
      @permanent_disability_status = params[:permanentDisabilityStatus]
      @pension_award_status_indicator = params[:pensionAwardStatusIndicator]
      @individual_ratings = params[:individualRatings] || []
      @metadata = params[:metadata]
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          vadrResponse: {
            type: :object,
            properties: {
              ptdrInformation: {
                type: :object,
                properties: {
                  serviceConnectedStatusIndicator: { type: :boolean, example: true },
                  pensionAwardStatusIndicator: { type: :boolean, example: false },
                  totalDisabilityStatusIndicator: { type: :boolean, example: true },
                  totalDisabilityEffectiveDate: { type: :string, example: "2023-01-01" },
                  responseMetadata: {
                    type: :object,
                    properties: {
                      responseCode: { type: :string, example: "MS000000" },
                      responseText: { type: :string, example: "Success" }
                    }
                  }
                }
              },
              sdrInformation: {
                type: :object,
                properties: {
                  combinedDisabilityRatingPercentage: { type: :integer, example: 100 },
                  combinedEffectiveDate: { type: :string, example: "2023-01-01" },
                  legalEffectiveDate: { type: :string, example: "2023-01-01" },
                  individualRatings: {
                    type: :array,
                    items: {
                      type: :object,
                      properties: {
                        decisionText: { type: :string, example: "Service Connection" },
                        ratingEffectiveDate: { type: :string, example: "2023-01-01" },
                        ratingEndDate: { type: :string, example: "2024-01-01" },
                        ratingPercentage: { type: :integer, example: 50 },
                        disabilityRatingId: { type: :string, example: "12345" },
                        staticIndicator: { type: :boolean, example: true }
                      }
                    }
                  },
                  responseMetadata: {
                    type: :object,
                    properties: {
                      responseCode: { type: :string, example: "MS000000" },
                      responseText: { type: :string, example: "Success" }
                    }
                  }
                }
              },
              earliestRatingEndDate: { type: :string, example: "2024-06-01" }
            }
          }
        }
      }
    end

    def as_json(options = {})
      {
        vadrResponse: {
          ptdrInformation: {
            serviceConnectedStatusIndicator: permanent_disability_status,
            pensionAwardStatusIndicator: pension_award_status_indicator,
            totalDisabilityStatusIndicator: total_disability_status,
            totalDisabilityEffectiveDate: total_disability_status_effective_date,
            responseMetadata: {
              responseCode: "MS000000",
              responseText: "Success"
            }
          },
          sdrInformation: {
            combinedDisabilityRatingPercentage: combined_disability_rating,
            combinedEffectiveDate: combined_effective_date,
            legalEffectiveDate: legal_effective_date,
            individualRatings: individual_ratings.map do |rating|
              {
                decisionText: rating[:decisionText],
                ratingEffectiveDate: rating[:ratingEffectiveDate],
                ratingEndDate: rating[:ratingEndDate],
                ratingPercentage: rating[:ratingPercentage],
                disabilityRatingId: rating[:disabilityRatingId],
                staticIndicator: rating[:staticIndicator]
              }
            end,
            responseMetadata: {
              responseCode: "MS000000",
              responseText: "Success"
            }
          },
          earliestRatingEndDate: earliest_rating_end_date
        }
      }
    end
  end
end
