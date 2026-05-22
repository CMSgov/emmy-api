module Api
  module V0
    class BatchEducationEnrollmentsController < ApplicationController
      def create
        reporter = Reporting::Reporter.new

        batch_params = Education::BatchStudentRequest.new(params)

        # Validate request
        if batch_params.batch_id.blank? || batch_params.submitted_by.blank? || batch_params.students.blank?
          return render json: { error: "missing required fields: batchId, submittedBy, and students are required" }, status: :bad_request
        end

        coordinator = Education::ServiceCoordinator.new
        begin
          coordinator.register_batch(batch_params)

          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "NSC",
            client_id: client_id,
            status_code: 201,
            success: true
          ))

          response_data = Education::BatchStudentCreatedResponse.new(
            message: "Batch registration successful",
            batch_job_id: batch_params.batch_id
          )
          render json: response_data, status: :created
        rescue StandardError => e
          Rails.logger.error("Batch education registration failed: #{e.message} #{e.backtrace.join("\n")}")
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "NSC",
            client_id: client_id,
            status_code: 500,
            success: false
          ))
          render json: { error: "Internal Server Error" }, status: :internal_server_error
        end
      end

      def show
        reporter = Reporting::Reporter.new
        coordinator = Education::ServiceCoordinator.new
        batch_request = Education::BatchStatusRequest.new(params)
        begin
          status_data = coordinator.get_batch_status(batch_request.id)
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "Postgres",
            client_id: client_id,
            status_code: 200,
            success: true
          ))
          response_data = Education::BatchStatusResponse.new(status_data)
          render json: response_data, status: :ok
        rescue ActiveRecord::RecordNotFound
          render json: { error: "Batch not found" }, status: :not_found
        rescue StandardError => e
          Rails.logger.error("Get batch status failed: #{e.message}")
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "Postgres",
            client_id: client_id,
            status_code: 500,
            success: false
          ))
          render json: { error: "Internal Server Error" }, status: :internal_server_error
        end
      end

      def details
        reporter = Reporting::Reporter.new
        coordinator = Education::ServiceCoordinator.new
        batch_request = Education::BatchDetailsRequest.new(params)
        begin
          details_data = coordinator.get_batch_details(batch_request.batch_job_id)
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "Postgres",
            client_id: client_id,
            status_code: 200,
            success: true
          ))
          response_data = Education::BatchDetailsResponse.new(details_data)
          render json: response_data, status: :ok
        rescue ActiveRecord::RecordNotFound
          render json: { error: "Batch not found" }, status: :not_found
        rescue StandardError => e
          Rails.logger.error("Get batch details failed: #{e.message}")
          reporter.report(Reporting::ReportData.new(
            timestamp: Time.now,
            endpoint: request.path,
            data_source: "Postgres",
            client_id: client_id,
            status_code: 500,
            success: false
          ))
          render json: { error: "Internal Server Error" }, status: :internal_server_error
        end
      end
    end
  end
end
