class AddUniqueIndexToEducationBatchStudentResults < ActiveRecord::Migration[8.1]
  def change
    remove_index :education_batch_student_results, name: "idx_on_education_batch_student_id_662efa214a"
    add_index :education_batch_student_results, :education_batch_student_id, unique: true
  end
end
