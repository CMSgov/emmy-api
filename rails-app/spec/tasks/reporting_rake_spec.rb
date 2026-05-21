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
      new-client,Dummy Agency
      missing-name,
    CSV
    csv.close

    expect {
      Rake::Task['reporting:import_client_agencies'].invoke(csv.path)
    }.to output(/created=1 updated=1 skipped=1/).to_stdout

    expect(ClientAgency.find('existing-client').agency_name).to eq('New Name')
    expect(ClientAgency.find('new-client').agency_name).to eq('Dummy Agency')
  ensure
    csv.unlink if csv
  end

  it 'accepts client_name as the CSV header for the display name' do
    csv = Tempfile.new(['client-agencies-alt-header', '.csv'])
    csv.write(<<~CSV)
      client_id,client_name
      named-client,Dummy Client
    CSV
    csv.close

    Rake::Task['reporting:import_client_agencies'].invoke(csv.path)

    expect(ClientAgency.find('named-client').agency_name).to eq('Dummy Client')
  ensure
    csv.unlink if csv
  end
end

RSpec.describe 'reporting:import_cognito_clients' do
  before(:all) do
    Rails.application.load_tasks if Rake::Task.tasks.empty?
  end

  before do
    Rake::Task['reporting:import_cognito_clients'].reenable
    allow(Kernel).to receive(:require).and_call_original
    allow(Kernel).to receive(:require).with('aws-sdk-cognitoidentityprovider').and_return(true)
  end

  it 'imports all paginated Cognito clients and updates existing mappings' do
    stub_const('Aws', Module.new)
    stub_const('Aws::CognitoIdentityProvider', Module.new)
    stub_const('Aws::CognitoIdentityProvider::Client', Class.new)

    first_page = instance_double(
      'ListUserPoolClientsResponse',
      user_pool_clients: [
        instance_double('UserPoolClientDescription', client_id: 'existing-client', client_name: 'New Name'),
        instance_double('UserPoolClientDescription', client_id: 'new-client', client_name: 'Dummy Client')
      ],
      next_token: 'page-2'
    )

    second_page = instance_double(
      'ListUserPoolClientsResponse',
      user_pool_clients: [
        instance_double('UserPoolClientDescription', client_id: 'skip-client', client_name: nil)
      ],
      next_token: nil
    )

    cognito_client = instance_double('Aws::CognitoIdentityProvider::Client')
    allow(Aws::CognitoIdentityProvider::Client).to receive(:new).with(region: 'us-east-1').and_return(cognito_client)
    allow(cognito_client).to receive(:list_user_pool_clients).with(
      user_pool_id: 'us-east-1_example',
      max_results: 60,
      next_token: nil
    ).and_return(first_page)
    allow(cognito_client).to receive(:list_user_pool_clients).with(
      user_pool_id: 'us-east-1_example',
      max_results: 60,
      next_token: 'page-2'
    ).and_return(second_page)

    ClientAgency.create!(client_id: 'existing-client', agency_name: 'Old Name')

    ClimateControl.modify(
      COGNITO_REGION: 'us-east-1',
      COGNITO_USER_POOL_ID: 'us-east-1_example'
    ) do
      expect {
        Rake::Task['reporting:import_cognito_clients'].invoke
      }.to output(/seen=3 created=1 updated=1 skipped=1/).to_stdout
    end

    expect(ClientAgency.find('existing-client').agency_name).to eq('New Name')
    expect(ClientAgency.find('new-client').agency_name).to eq('Dummy Client')
  end
end
