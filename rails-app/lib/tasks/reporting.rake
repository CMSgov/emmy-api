require 'csv'

namespace :reporting do
  desc 'Import client agency mappings from a CSV file with client_id and agency_name/client_name columns'
  task :import_client_agencies, [:path] => :environment do |_task, args|
    path = args[:path].presence || ENV['FILE']
    abort 'Usage: bin/rake reporting:import_client_agencies[path/to/file.csv] or FILE=path/to/file.csv' if path.blank?
    abort "CSV file not found: #{path}" unless File.exist?(path)

    created = 0
    updated = 0
    skipped = 0

    CSV.foreach(path, headers: true) do |row|
      client_id = row['client_id']&.strip
      agency_name = row['agency_name']&.strip || row['client_name']&.strip || row['name']&.strip

      if client_id.blank? || agency_name.blank?
        skipped += 1
        next
      end

      record = ClientAgency.find_or_initialize_by(client_id: client_id)

      if record.new_record?
        record.agency_name = agency_name
        record.updated_at = Time.current
        record.save!
        created += 1
        next
      end

      if record.agency_name == agency_name
        skipped += 1
        next
      end

      record.update!(agency_name: agency_name, updated_at: Time.current)
      updated += 1
    end

    puts "Imported client agencies from #{path}: created=#{created} updated=#{updated} skipped=#{skipped}" # rubocop:disable Rails/Output
  end

  desc 'Import Cognito app clients into client_agencies using COGNITO_REGION and COGNITO_USER_POOL_ID'
  task import_cognito_clients: :environment do
    require 'aws-sdk-cognitoidentityprovider'

    region = ENV.fetch('COGNITO_REGION')
    user_pool_id = ENV.fetch('COGNITO_USER_POOL_ID')

    cognito = Aws::CognitoIdentityProvider::Client.new(region: region)
    created = 0
    updated = 0
    skipped = 0
    seen = 0
    next_token = nil

    loop do
      response = cognito.list_user_pool_clients(
        user_pool_id: user_pool_id,
        max_results: 60,
        next_token: next_token
      )

      response.user_pool_clients.each do |app_client|
        seen += 1
        client_id = app_client.client_id&.strip
        agency_name = app_client.client_name&.strip

        if client_id.blank? || agency_name.blank?
          skipped += 1
          next
        end

        record = ClientAgency.find_or_initialize_by(client_id: client_id)

        if record.new_record?
          record.agency_name = agency_name
          record.updated_at = Time.current
          record.save!
          created += 1
          next
        end

        if record.agency_name == agency_name
          skipped += 1
          next
        end

        record.update!(agency_name: agency_name, updated_at: Time.current)
        updated += 1
      end

      next_token = response.next_token
      break if next_token.blank?
    end

    puts "Imported Cognito clients: seen=#{seen} created=#{created} updated=#{updated} skipped=#{skipped}" # rubocop:disable Rails/Output
  end
end
