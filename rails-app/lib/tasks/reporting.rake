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
end
