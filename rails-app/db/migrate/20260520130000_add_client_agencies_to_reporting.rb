class AddClientAgenciesToReporting < ActiveRecord::Migration[8.1]
  def change
    add_column :api_events, :agency_name, :string unless column_exists?(:api_events, :agency_name)
    add_index :api_events, :agency_name unless index_exists?(:api_events, :agency_name)

    return if table_exists?(:client_agencies)

    create_table :client_agencies, id: false do |t|
      t.string :client_id, null: false
      t.string :agency_name, null: false
      t.datetime :updated_at, null: false, default: -> { 'CURRENT_TIMESTAMP' }
    end

    add_index :client_agencies, :client_id, unique: true
  end
end
