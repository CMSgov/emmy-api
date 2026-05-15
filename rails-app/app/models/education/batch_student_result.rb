# == Schema Information
#
# Table name: education_batch_student_results
#
#  id                         :uuid             not null, primary key
#  found_enrollment           :boolean
#  results                    :jsonb
#  created_at                 :datetime         not null
#  updated_at                 :datetime         not null
#  education_batch_student_id :uuid             not null
#
# Indexes
#
#  idx_on_education_batch_student_id_662efa214a  (education_batch_student_id)
#
# Foreign Keys
#
#  fk_rails_...  (education_batch_student_id => education_batch_students.id)
#
class Education::BatchStudentResult < ApplicationRecord
  belongs_to :education_batch_student, class_name: "Education::BatchStudent"
end
