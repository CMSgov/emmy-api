# == Schema Information
#
# Table name: education_enrollment_batches
#
#  id           :uuid             not null, primary key
#  callback_url :string           not null
#  status       :string           not null
#  submitted_by :string           not null
#  created_at   :datetime         not null
#  updated_at   :datetime         not null
#  batch_id     :string           not null
#
# Indexes
#
#  index_education_enrollment_batches_on_batch_id  (batch_id) UNIQUE
#
class Education::EnrollmentBatch < ApplicationRecord
  has_many :education_batch_students, class_name: "Education::BatchStudent", foreign_key: :education_enrollment_batch_id, dependent: :destroy

  validates :batch_id, presence: true, uniqueness: true
  validates :submitted_by, presence: true
  validates :callback_url, presence: true
  enum :status, {
    queued: 'QUEUED',
    processing: 'PROCESSING',
    completed: 'COMPLETED',
    failed: 'FAILED'
  }, default: 'QUEUED', prefix: true
end
