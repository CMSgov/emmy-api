module Api
  module V0
    class EducationEnrollmentsController < ApplicationController
      def create
        reporter = Reporting::Reporter.new
        client_id = request.headers["Authorization"] # Simplified placeholder if no better option exists

        if missing_field = missing_education_identity_field
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "NSC",
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

        coordinator = Education::ServiceCoordinator.new
        begin
          result = coordinator.lookup_enrollment_status(req_params)
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "NSC",
            client_id: client_id,
            status_code: 200,
            success: true
          ))
          render json: result, status: :ok
        rescue Education::NotFoundError => e
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "NSC",
            client_id: client_id,
            status_code: 404,
            success: false
          ))
          render json: { error: e.message }, status: :not_found
        rescue StandardError => e
          # Log the error like Go does
          Rails.logger.error("Education verification failed: #{e.message} #{e.backtrace.join("\n")}")
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "NSC",
            client_id: client_id,
            status_code: 502,
            success: false
          ))
          render json: { error: "Bad Gateway" }, status: :bad_gateway
        end
      end

      private

      def missing_education_identity_field
        return "firstName" if params[:firstName].blank?
        return "lastName" if params[:lastName].blank?
        return "dateOfBirth" if params[:dateOfBirth].blank?
        nil
      end
    end
  end
end
