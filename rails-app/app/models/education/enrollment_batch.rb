class Education::EnrollmentBatch < ApplicationRecord
  has_many :education_batch_students, class_name: 'Education::BatchStudent', foreign_key: :education_enrollment_batch_id, dependent: :destroy

  validates :batch_id, presence: true, uniqueness: true
  validates :submitted_by, presence: true
  validates :callback_url, presence: true
  validates :status, presence: true
end
