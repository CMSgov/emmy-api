module Education
  class ServiceCoordinator
    def initialize
      @client = Education::NscClient.new
      @sqs_client = Aws::SQS::Client.new if ENV["BATCH_SQS_QUEUE_URL"].present?
    end

    def lookup_enrollment_status(req_params)
      light = Stoplight("nsc-lookup-enrollment-status")

      # We don't want to trip the circuit for NotFoundError
      if light.respond_to?(:with_allowed_errors)
        light = light.with_allowed_errors([ Education::NotFoundError ])
      end

      light.run do
        @client.lookup_enrollment_status(req_params)
      end
    end

    def register_batch(batch_student_request)
      batch = ActiveRecord::Base.transaction do
        b = Education::EnrollmentBatch.create!(
          batch_id: batch_student_request.batch_id,
          submitted_by: batch_student_request.submitted_by,
          callback_url: batch_student_request.callback_url,
          status: :queued
        )

        batch_student_request.students.each do |student_params|
          b.education_batch_students.create!(
            record_id: student_params[:record_id],
            first_name: student_params[:first_name],
            last_name: student_params[:last_name],
            date_of_birth: student_params[:date_of_birth],
            ssn: student_params[:ssn],
            status: :queued
          )
        end
        b
      end

      enqueue_batch_job(batch.id) if batch
      batch
    end

    def get_batch_status(batch_id)
      batch = Education::EnrollmentBatch.find_by!(batch_id: batch_id)

      # Using a manual count approach similar to the Go implementation
      students = batch.education_batch_students
      total_records = students.count
      processed_records = students.where(status: [ :success, :failed, :no_hit ]).count
      success_count = students.where(status: :success).count
      failure_count = students.where(status: :failed).count

      {
        batch_job_id: batch.batch_id,
        status: batch.status.upcase,
        submitted_at: batch.created_at,
        updated_at: Time.current,
        total_records: total_records,
        processed_records: processed_records,
        success_count: success_count,
        failure_count: failure_count
      }
    end

    def get_batch_details(batch_id)
      batch = Education::EnrollmentBatch.find_by!(batch_id: batch_id)

      results = batch.education_batch_students.includes(:education_batch_student_result).map do |student|
        result = student.education_batch_student_result
        {
          record_id: student.record_id,
          status: student.status.upcase,
          found_enrollment: result&.found_enrollment || false,
          results: result&.results
        }
      end

      {
        batch_job_id: batch.batch_id,
        results: results
      }
    end

    private

    def enqueue_batch_job(enrollment_batch_id)
      queue_url = ENV["BATCH_SQS_QUEUE_URL"]
      return unless queue_url.present? && @sqs_client

      @sqs_client.send_message(
        queue_url: queue_url,
        message_body: { enrollment_batch_id: enrollment_batch_id }.to_json
      )
    rescue StandardError => e
      Rails.logger.error("Failed to enqueue batch job to SQS: #{e.message}")
    end
  end
end
