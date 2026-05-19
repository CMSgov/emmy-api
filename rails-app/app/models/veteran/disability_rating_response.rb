module Veteran
  class DisabilityRatingResponse
    attr_accessor :raw_data, :data_source, :metadata,
                  :total_disability_status, :total_disability_status_effective_date,
                  :combined_disability_rating, :combined_effective_date, :legal_effective_date,
                  :earliest_rating_end_date

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
    end

    def as_json(options = {})
      {
        totalDisabilityStatus: total_disability_status,
        totalDisabilityStatusEffectiveDate: total_disability_status_effective_date,
        combinedDisabilityRating: combined_disability_rating,
        combinedEffectiveDate: combined_effective_date,
        legalEffectiveDate: legal_effective_date,
        earliestRatingEndDate: earliest_rating_end_date,
        dataSource: data_source,
        metadata: metadata,
        rawData: raw_data
      }
    end
  end
end
