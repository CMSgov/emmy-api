module Api
  module V0
    class EducationEnrollmentsController < ApplicationController
      def create
        if missing_field = missing_education_identity_field
          return render json: { error: "missing required field: #{missing_field}" }, status: :bad_request
        end

        req_params = params.permit(
          :firstName, :lastName, :dateOfBirth, :ssn, :middleName,
          address: [:street1, :street2, :street3, :city, :state, :postalCode, :country]
        ).to_h.deep_symbolize_keys

        coordinator = Education::ServiceCoordinator.new
        begin
          result = coordinator.lookup_enrollment_status(req_params)
          render json: result, status: :ok
        rescue Education::NotFoundError => e
          render json: { error: e.message }, status: :not_found
        rescue StandardError => e
          # Log the error like Go does
          Rails.logger.error("Education verification failed: #{e.message} #{e.backtrace.join("\n")}")
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
