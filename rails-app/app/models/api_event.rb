# == Schema Information
#
# Table name: api_events
#
#  id          :integer          not null, primary key
#  data_source :text             not null
#  endpoint    :text             not null
#  status_code :integer          not null
#  success     :boolean          not null
#  timestamp   :timestamptz      not null
#  created_at  :timestamptz      not null
#  client_id   :text             not null
#
# Indexes
#
#  idx_api_events_client_id       (client_id)
#  idx_api_events_timestamp       (timestamp)
#  index_api_events_on_client_id  (client_id)
#  index_api_events_on_timestamp  (timestamp)
#
class ApiEvent < ApplicationRecord
end
