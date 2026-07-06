module Veteran
  class DisabilityRatingResponseV0
    attr_accessor :raw_data, :data_source, :metadata,
                  :total_disability_status, :total_disability_status_effective_date,
                  :combined_disability_rating, :combined_effective_date, :legal_effective_date,
                  :earliest_rating_end_date, :permanent_disability_status

    def initialize(params = {})
      @raw_data = params[:rawData]
      @data_source = params[:dataSource]
      @metadata = params[:metadata]
      @total_disability_status = params[:totalDisabilityStatus]
      @total_disability_status_effective_date = params[:totalDisabilityStatusEffectiveDate]
      @combined_disability_rating = params[:combinedDisabilityRating]
      @combined_effective_date = params[:combinedEffectiveDate]
      @legal_effective_date = params[:legalEffectiveDate]
      @earliest_rating_end_date = params[:earliestRatingEndDate]
      @permanent_disability_status = params[:permanentDisabilityStatus]
    end

    def self.to_swagger_schema
      {
        type: :object,
        properties: {
          permanentDisabilityStatus: { type: :boolean, description: "Whether they have been awarded permanent disability status for work purposes", example: true },
          totalDisabilityStatus: { type: :boolean, description: "Whether the veteran is totally disabled.", example: true },
          totalDisabilityStatusEffectiveDate: { type: :string, description: "Effective date of total disability status (YYYY-MM-DD).", example: "2023-01-01" },
          combinedDisabilityRating: { type: :integer, description: "Combined disability rating percentage.", example: 100 },
          combinedEffectiveDate: { type: :string, description: "Effective date of combined disability rating (YYYY-MM-DD).", example: "2023-01-01" },
          legalEffectiveDate: { type: :string, description: "Legal effective date of disability rating (YYYY-MM-DD).", example: "2023-01-01" },
          earliestRatingEndDate: { type: :string, description: "Earliest end date for any individual disability rating (YYYY-MM-DD).", example: "2024-06-01" },
          dataSource: { type: :string, description: "The source of the disability rating data (e.g., VA).", example: "VA" },
          metadata: {
            type: :object,
            properties: {
              apiVersion: { type: :string, example: "1.3.0" },
              environment: { type: :string, example: "development" },
              requestTimestamp: { type: :string, example: "2023-11-01T12:00:00.000Z" },
              responseTimestamp: { type: :string, example: "2023-11-01T12:00:00.100Z" },
              datasourceDurationMillis: { type: :integer, example: 100 },
              transactionId: { type: :string, example: "uuid" }
            }
          }
        }
      }
    end

    def as_json(options = {})
      {
        totalDisabilityStatus: total_disability_status,
        totalDisabilityStatusEffectiveDate: total_disability_status_effective_date,
        combinedDisabilityRating: combined_disability_rating,
        combinedEffectiveDate: combined_effective_date,
        legalEffectiveDate: legal_effective_date,
        earliestRatingEndDate: earliest_rating_end_date,
        permanentDisabilityStatus: permanent_disability_status,
        dataSource: data_source,
        metadata: metadata,
        rawData: raw_data
      }
    end
  end
end
