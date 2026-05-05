module Education
  class EnrollmentResponse
    attr_accessor :enrollment_status, :enrollment_details, :raw_data, :data_source, :metadata

    def initialize(params = {})
      @enrollment_status = params[:enrollmentStatus]
      @enrollment_details = params[:enrollmentDetails] || []
      @raw_data = params[:rawData]
      @data_source = params[:dataSource]
      @metadata = params[:metadata]
    end

    def as_json(options = {})
      {
        enrollmentStatus: enrollment_status,
        enrollmentDetails: enrollment_details,
        rawData: raw_data,
        dataSource: data_source,
        metadata: metadata
      }
    end
  end
end
