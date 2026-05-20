require 'rails_helper'
require 'rake'
require 'tempfile'

RSpec.describe 'reporting:import_client_agencies' do
  before(:all) do
    Rails.application.load_tasks if Rake::Task.tasks.empty?
  end

  before do
    Rake::Task['reporting:import_client_agencies'].reenable
  end

  it 'imports new mappings and updates existing ones from CSV' do
    ClientAgency.create!(client_id: 'existing-client', agency_name: 'Old Name')

    csv = Tempfile.new(['client-agencies', '.csv'])
    csv.write(<<~CSV)
      client_id,agency_name
      existing-client,New Name
      new-client,CMS/DSAC
      missing-name,
    CSV
    csv.close

    expect {
      Rake::Task['reporting:import_client_agencies'].invoke(csv.path)
    }.to output(/created=1 updated=1 skipped=1/).to_stdout

    expect(ClientAgency.find('existing-client').agency_name).to eq('New Name')
    expect(ClientAgency.find('new-client').agency_name).to eq('CMS/DSAC')
  ensure
    csv.unlink if csv
  end

  it 'accepts client_name as the CSV header for the display name' do
    csv = Tempfile.new(['client-agencies-alt-header', '.csv'])
    csv.write(<<~CSV)
      client_id,client_name
      named-client,Imhotep
    CSV
    csv.close

    Rake::Task['reporting:import_client_agencies'].invoke(csv.path)

    expect(ClientAgency.find('named-client').agency_name).to eq('Imhotep')
  ensure
    csv.unlink if csv
  end
end
