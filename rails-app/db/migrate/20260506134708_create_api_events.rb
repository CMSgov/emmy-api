class CreateApiEvents < ActiveRecord::Migration[8.1]
  def change
    create_table :api_events, if_not_exists: true do |t|
      t.datetime :timestamp
      t.string :endpoint
      t.string :data_source
      t.string :client_id
      t.integer :status_code
      t.boolean :success

      t.timestamps
    end
    # add_index :api_events, :timestamp
    # add_index :api_events, :client_id
  end
end
