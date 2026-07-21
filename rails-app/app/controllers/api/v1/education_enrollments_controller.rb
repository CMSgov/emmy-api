module Api
  module V1
    class EducationEnrollmentsController < ApplicationController
      def create
        reporter = Reporting::Reporter.new

        if params[:nscRequest].blank?
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "NSC",
            client_id: client_id,
            status_code: 400,
            success: false
          ))
          return render json: { error: "missing required nscRequest object" }, status: :bad_request
        end

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

        req_params = params.require(:nscRequest).permit(
          :personSocialSecurityNumber, :personGivenName, :personMiddleName, :personSurName, :personBirthDate, :asOfDate, :termsAcceptedIndicator,
          previousNames: [ :personGivenName, :personMiddleName, :personSurName ]
        ).to_h.deep_symbolize_keys
        enrollment_req = Education::EnrollmentRequestV1.new(req_params)

        coordinator = Education::ServiceCoordinator.new
        begin
          result = coordinator.lookup_enrollment_status(enrollment_req)
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
        nsc_req = params[:nscRequest] || {}
        return nil if nsc_req[:personSocialSecurityNumber].present?

        return "personGivenName" if nsc_req[:personGivenName].blank?
        return "personSurName" if nsc_req[:personSurName].blank?
        return "personBirthDate" if nsc_req[:personBirthDate].blank?
        nil
      end
    end
  end
end
