require 'rails_helper'

module Education
  RSpec.describe BatchStudentResult, type: :model do
    let(:batch) do
      Education::EnrollmentBatch.create!(
        batch_id: 'test-batch-unique',
        submitted_by: 'test-user',
        callback_url: 'http://example.com/callback',
        status: :queued
      )
    end

    let(:student) do
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

    it 'enforces uniqueness on education_batch_student_id' do
      Education::BatchStudentResult.create!(
        education_batch_student_id: student.id,
        found_enrollment: true,
        results: { status: 'found' }
      )

      duplicate = Education::BatchStudentResult.new(
        education_batch_student_id: student.id,
        found_enrollment: false,
        results: { status: 'not found' }
      )

      expect(duplicate).not_to be_valid
      expect(duplicate.errors[:education_batch_student_id]).to include('has already been taken')

      expect {
        duplicate.save!(validate: false)
      }.to raise_error(ActiveRecord::RecordNotUnique)
    end
  end
end
