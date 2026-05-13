module Veteran
  class DisabilityRatingResponse
    attr_accessor :raw_data, :data_source, :metadata,
                  :total_disability_status, :total_disability_status_effective_date

    def initialize(params = {})
      @raw_data = params[:rawData]
      @data_source = params[:dataSource]
      @metadata = params[:metadata]
      @total_disability_status = params[:totalDisabilityStatus]
      @total_disability_status_effective_date = params[:totalDisabilityStatusEffectiveDate]
      @permanent_and_total_disability_status = params[:permanentAndTotalDisabilityStatus]
      @permanent_and_total_disability_pension_award_status = params[:permanentAndTotalDisabilityPensionAwardStatus]
    end

    def as_json(options = {})
      {
        totalDisabilityStatus: total_disability_status,
        totalDisabilityStatusEffectiveDate: total_disability_status_effective_date,
        dataSource: data_source,
        metadata: metadata,
        rawData: raw_data
      }
    end
  end
end
