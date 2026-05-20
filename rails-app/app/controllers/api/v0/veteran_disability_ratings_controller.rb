module Api
  module V0
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

        req_params = params.permit(
          :firstName, :lastName, :dateOfBirth, :ssn, :middleName,
          address: [ :street1, :street2, :street3, :city, :state, :postalCode, :country ]
        ).to_h.deep_symbolize_keys

        coordinator = Veteran::ServiceCoordinator.new
        begin
          result = coordinator.lookup_disability_rating(req_params)
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
          Rails.logger.error("Veteran disability lookup failed: #{e.message} #{e.backtrace&.join("\n")}")
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
        return "firstName" if params[:firstName].blank?
        return "lastName" if params[:lastName].blank?
        return "dateOfBirth" if params[:dateOfBirth].blank?

        # In Go, it checks for SSN OR complete address
        if params[:ssn].blank? && !complete_address?
          return "ssn or a complete address"
        end

        nil
      end

      def complete_address?
        addr = params[:address]
        return false if addr.blank?

        addr[:street1].present? &&
          addr[:city].present? &&
          addr[:state].present? &&
          addr[:postalCode].present? &&
          addr[:country].present?
      end
    end
  end
end
