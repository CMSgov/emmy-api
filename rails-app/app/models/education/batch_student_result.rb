class Education::BatchStudentResult < ApplicationRecord
  belongs_to :education_batch_student, class_name: 'Education::BatchStudent'
end
