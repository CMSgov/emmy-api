module Veteran
  class DisabilityRatingResponse
    attr_accessor :combined_disability_rating, :combined_effective_date, :legal_effective_date,
                  :earliest_rating_end_date, :raw_data, :data_source, :metadata

    def initialize(params = {})
      @combined_disability_rating = params[:combinedDisabilityRating]
      @combined_effective_date = params[:combinedEffectiveDate]
      @legal_effective_date = params[:legalEffectiveDate]
      @earliest_rating_end_date = params[:earliestRatingEndDate]
      @raw_data = params[:rawData]
      @data_source = params[:dataSource]
      @metadata = params[:metadata]
    end

    def as_json(options = {})
      {
        combinedDisabilityRating: combined_disability_rating,
        combinedEffectiveDate: combined_effective_date,
        legalEffectiveDate: legal_effective_date,
        earliestRatingEndDate: earliest_rating_end_date,
        rawData: raw_data,
        dataSource: data_source,
        metadata: metadata
      }
    end
  end
end
