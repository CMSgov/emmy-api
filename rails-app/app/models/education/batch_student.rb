class Education::BatchStudent < ApplicationRecord
  belongs_to :education_enrollment_batch, class_name: 'Education::EnrollmentBatch'
  has_one :education_batch_student_result, class_name: 'Education::BatchStudentResult', foreign_key: :education_batch_student_id, dependent: :destroy

  encrypts :ssn

  validates :record_id, presence: true
  validates :first_name, presence: true
  validates :last_name, presence: true
  validates :date_of_birth, presence: true
  validates :status, presence: true
end
