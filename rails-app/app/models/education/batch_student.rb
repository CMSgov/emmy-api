# == Schema Information
#
# Table name: education_batch_students
#
#  id                            :uuid             not null, primary key
#  date_of_birth                 :string           not null
#  first_name                    :string           not null
#  last_name                     :string           not null
#  ssn                           :string           not null
#  status                        :string           not null
#  created_at                    :datetime         not null
#  updated_at                    :datetime         not null
#  education_enrollment_batch_id :uuid             not null
#  record_id                     :string           not null
#
# Indexes
#
#  idx_on_education_enrollment_batch_id_787e79cd79  (education_enrollment_batch_id)
#
# Foreign Keys
#
#  fk_rails_...  (education_enrollment_batch_id => education_enrollment_batches.id)
#
class Education::BatchStudent < ApplicationRecord
  belongs_to :education_enrollment_batch, class_name: "Education::EnrollmentBatch"
  has_one :education_batch_student_result, class_name: "Education::BatchStudentResult", foreign_key: :education_batch_student_id, dependent: :destroy

  encrypts :ssn

  validates :record_id, presence: true
  validates :first_name, presence: true
  validates :last_name, presence: true
  validates :date_of_birth, presence: true
  enum :status, {
    queued: "QUEUED",
    processing: "PROCESSING",
    success: "SUCCESS",
    no_hit: "NO_HIT",
    failed: "FAILED"
  }, default: "QUEUED", prefix: true
end
