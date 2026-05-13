class CreateEducationBatchStudentResults < ActiveRecord::Migration[8.1]
  def change
    create_table :education_batch_student_results, if_not_exists: true, id: :uuid do |t|
      t.references :education_batch_student, null: false, foreign_key: true, type: :uuid
      t.jsonb :results
      t.boolean :found_enrollment

      t.timestamps
    end
  end
end
