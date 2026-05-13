class CreateEducationBatchStudents < ActiveRecord::Migration[8.1]
  def change
    create_table :education_batch_students, if_not_exists: true, id: :uuid do |t|
      t.references :education_enrollment_batch, null: false, foreign_key: true, type: :uuid
      t.string :record_id, null: false
      t.string :first_name, null: false
      t.string :last_name, null: false
      t.string :date_of_birth, null: false
      t.string :ssn, null: false
      t.string :status, null: false

      t.timestamps
    end
  end
end
