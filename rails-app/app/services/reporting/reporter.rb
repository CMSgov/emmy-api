require "aws-sdk-sqs"

module Reporting
  class Reporter
    def initialize
      @queue_url = ENV["SQS_QUEUE_URL"]

      if @queue_url.present?
        @sqs_client = Aws::SQS::Client.new
      end
    end

    def report(data)
      if @queue_url.present? && @sqs_client
        begin
          @sqs_client.send_message(
            queue_url: @queue_url,
            message_body: data.to_json
          )
          Rails.logger.debug("report sent to SQS: #{@queue_url}")
        rescue StandardError => e
          Rails.logger.error("failed to send report to SQS: #{e.message}")
        end
      else
        Rails.logger.info(
          "api call report: endpoint=#{data.endpoint} success=#{data.success} " \
          "data_source=#{data.data_source} client_id=#{data.client_id} " \
          "timestamp=#{data.timestamp} status_code=#{data.status_code}"
        )
      end
    end
  end
end
