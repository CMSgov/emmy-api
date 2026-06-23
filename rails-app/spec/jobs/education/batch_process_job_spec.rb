require 'rails_helper'

module Education
  RSpec.describe BatchProcessJob, type: :job do
    let(:batch) do
      Education::EnrollmentBatch.create!(
        batch_id: 'test-batch-123',
        submitted_by: 'test-user',
        callback_url: 'http://example.com/callback',
        status: :queued
      )
    end

    let!(:student1) do
      Education::BatchStudent.create!(
        education_enrollment_batch: batch,
        record_id: 'rec1',
        first_name: 'John',
        last_name: 'Doe',
        date_of_birth: '1990-01-01',
        ssn: '123456789',
        status: :queued
      )
    end

    let!(:student2) do
      Education::BatchStudent.create!(
        education_enrollment_batch: batch,
        record_id: 'rec2',
        first_name: 'Jane',
        last_name: 'Smith',
        date_of_birth: '1992-02-02',
        ssn: '987654321',
        status: :queued
      )
    end

    let(:coordinator) { instance_double(Education::ServiceCoordinator) }
    let(:body) { { enrollment_batch_id: batch.id }.to_json }

    before do
      allow(Education::ServiceCoordinator).to receive(:new).and_return(coordinator)
    end

    around do |example|
      ActiveRecord::Encryption.config.primary_key = "test-primary-key"
      ActiveRecord::Encryption.config.deterministic_key = "test-deterministic-key"
      ActiveRecord::Encryption.config.key_derivation_salt = "test-salt"
      example.run
    ensure
      ActiveRecord::Encryption.config.primary_key = nil
      ActiveRecord::Encryption.config.deterministic_key = nil
      ActiveRecord::Encryption.config.key_derivation_salt = nil
    end

    describe '#perform' do
      it 'processes all students and updates batch status' do
        resp1 = Education::EnrollmentResponseV0.new(enrollmentStatus: 'FULL_TIME', dataSource: 'NSC')

        expect(coordinator).to receive(:lookup_enrollment_status).with(kind_of(Education::EnrollmentRequestV0)).twice do |req|
          if req.first_name == 'John'
            resp1
          elsif req.first_name == 'Jane'
            raise ::Education::NotFoundError
          else
            raise "Unexpected student: #{req.first_name}"
          end
        end

        BatchProcessJob.new.perform(nil, body)

        expect(batch.reload.status).to eq('completed')

        expect(student1.reload.status).to eq('success')
        expect(student1.education_batch_student_result).to be_present
        expect(student1.education_batch_student_result.found_enrollment).to be true

        expect(student2.reload.status).to eq('no_hit')
        expect(student2.education_batch_student_result).to be_present
        expect(student2.education_batch_student_result.found_enrollment).to be false
      end

      it 'handles unexpected errors for individual students' do
        allow(coordinator).to receive(:lookup_enrollment_status).and_raise(StandardError, "Unexpected error")

        BatchProcessJob.new.perform(nil, body)

        expect(batch.reload.status).to eq('completed')
        expect(student1.reload.status).to eq('failed')
        expect(student2.reload.status).to eq('failed')
      end

      it 'updates batch to FAILED if a major error occurs' do
        allow_any_instance_of(Education::EnrollmentBatch).to receive(:education_batch_students).and_raise(StandardError, "Database error")

        expect {
          BatchProcessJob.new.perform(nil, body)
        }.to raise_error(StandardError, "Database error")

        expect(batch.reload.status).to eq('failed')
      end

      it 'is idempotent and only processes queued or failed students' do
        student1.status_success!
        student2.status_failed!

        resp = Education::EnrollmentResponseV0.new(enrollmentStatus: 'FULL_TIME', dataSource: 'NSC')

        # Should only call for student2
        expect(coordinator).to receive(:lookup_enrollment_status).once do |req|
          expect(req.first_name).to eq('Jane')
          resp
        end

        BatchProcessJob.new.perform(nil, body)

        expect(batch.reload.status).to eq('completed')
        expect(student1.reload.status).to eq('success')
        expect(student2.reload.status).to eq('success')
      end
    end
  end
end
