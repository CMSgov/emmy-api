module Education
  class BatchProcessJob
    include Shoryuken::Worker

    # Ensure NotFoundError is loaded if it's defined in NscClient
    require_dependency "education/nsc_client"

    shoryuken_options queue: ENV["BATCH_SQS_QUEUE_URL"] || "education-batch-processing", auto_delete: true

    def perform(_sqs_msg, body)
      data = JSON.parse(body)
      batch_id = data["enrollment_batch_id"]

      batch = Education::EnrollmentBatch.find_by(id: batch_id)
      return unless batch

      batch.status_processing!

      coordinator = Education::ServiceCoordinator.new

      batch.education_batch_students.where(status: [ :queued, :failed ]).find_each do |student|
        process_student(coordinator, student)
      end

      batch.status_completed!
    rescue => e
      Rails.logger.error("Education::BatchProcessJob failed for batch #{batch_id}: #{e.message}")
      batch&.status_failed!
      raise e
    end

    private

    def process_student(coordinator, student)
      student.status_processing!

      params = {
        firstName: student.first_name,
        lastName: student.last_name,
        dateOfBirth: student.date_of_birth.to_s,
        ssn: student.ssn
      }

      begin
        response = coordinator.lookup_enrollment_status(params)

        Education::BatchStudentResult.create_or_find_by!(
          education_batch_student_id: student.id
        ).update!(
          found_enrollment: true,
          results: response.as_json
        )

        student.status_success!
      rescue ::Education::NotFoundError
        Education::BatchStudentResult.create_or_find_by!(
          education_batch_student_id: student.id
        ).update!(
          found_enrollment: false,
          results: { message: "No enrollment found" }
        )
        student.status_no_hit!
      rescue => e
        Rails.logger.error("Failed to process student #{student.id}: #{e.message}")
        student.status_failed!
        # We don't raise here to allow other students in the batch to be processed
      end
    end
  end
end
