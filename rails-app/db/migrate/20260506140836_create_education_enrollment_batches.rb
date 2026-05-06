class CreateEducationEnrollmentBatches < ActiveRecord::Migration[8.1]
  def change
    create_table :education_enrollment_batches, if_not_exists: true, id: :uuid do |t|
      t.string :batch_id, null: false
      t.string :submitted_by, null: false
      t.string :callback_url, null: false
      t.string :status, null: false

      t.timestamps
    end
    # add_index :education_enrollment_batches, :batch_id, unique: true
  end
end
