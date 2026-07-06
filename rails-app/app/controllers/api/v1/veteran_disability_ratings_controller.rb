module Api
  module V1
    class VeteranDisabilityRatingsController < ApplicationController
      def create
        reporter = Reporting::Reporter.new

        if missing_field = missing_veteran_identity_field
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "VA",
            client_id: client_id,
            status_code: 400,
            success: false
          ))
          return render json: { error: "missing required field: #{missing_field}" }, status: :bad_request
        end

        req_v1 = Veteran::DisabilityRatingRequestV1.new(params)
        coordinator = Veteran::ServiceCoordinator.new

        begin
          result = coordinator.lookup_disability_rating(req_v1, "V1")

          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "VA",
            client_id: client_id,
            status_code: 200,
            success: true
          ))
          render json: result, status: :ok
        rescue Veteran::NotFoundError => e
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "VA",
            client_id: client_id,
            status_code: 404,
            success: false
          ))
          render json: { error: e.message }, status: :not_found
        rescue StandardError => e
          Rails.logger.error("Veteran disability lookup V1 failed: #{e.message} #{e.backtrace&.join("\n")}")
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "VA",
            client_id: client_id,
            status_code: 502,
            success: false
          ))
          render json: { error: "Bad Gateway" }, status: :bad_gateway
        end
      end

      private

      def missing_veteran_identity_field
        vadr = params[:vadrRequest]
        return "vadrRequest" if vadr.blank?
        return "personGivenName" if vadr[:personGivenName].blank?
        return "personSurName" if vadr[:personSurName].blank?
        return "personBirthDate" if vadr[:personBirthDate].blank?

        nil
      end
    end
  end
end
