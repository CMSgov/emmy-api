class ClientAgency < ApplicationRecord
  self.primary_key = :client_id

  validates :client_id, presence: true, uniqueness: true
  validates :agency_name, presence: true
end
