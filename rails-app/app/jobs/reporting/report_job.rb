module Reporting
  class ReportJob
    include Shoryuken::Worker

    # The queue name should ideally be configured via environment variable
    shoryuken_options queue: ENV['SQS_QUEUE_NAME'] || 'reporting-events', auto_delete: true

    def perform(_sqs_msg, body)
      data = JSON.parse(body)

      Rails.logger.info("received report data: endpoint=#{data['endpoint']} success=#{data['success']} " \
                         "data_source=#{data['data_source']} client_id=#{data['client_id']} " \
                         "timestamp=#{data['timestamp']} status_code=#{data['status_code']}")

      ApiEvent.create!(
        timestamp: data['timestamp'],
        endpoint: data['endpoint'],
        data_source: data['data_source'],
        client_id: data['client_id'],
        status_code: data['status_code'],
        success: data['success']
      )
    rescue => e
      Rails.logger.error("Reporting::ReportJob failed to process message: #{e.message}")
      # In a real scenario, we might want to let it retry depending on the error
      raise e
    end
  end
end
