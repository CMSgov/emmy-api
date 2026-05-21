class AllowNullSsnOnBatchStudents < ActiveRecord::Migration[8.1]
  def change
    change_column_null :batch_students, :ssn, true
    change_column_null :education_batch_students, :ssn, true
  end
end
