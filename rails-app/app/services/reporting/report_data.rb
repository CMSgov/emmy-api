module Reporting
  class ReportData
    attr_accessor :timestamp, :endpoint, :data_source, :client_id, :status_code, :success, :transaction_id

    def initialize(timestamp:, endpoint:, data_source:, client_id:, status_code:, success:, transaction_id: Current.request_id)
      @timestamp = timestamp
      @endpoint = endpoint
      @data_source = data_source
      @client_id = client_id
      @status_code = status_code
      @success = success
      @transaction_id = transaction_id
    end

    def to_json(*_args)
      {
        timestamp: @timestamp.iso8601,
        endpoint: @endpoint,
        data_source: @data_source,
        client_id: @client_id,
        status_code: @status_code,
        success: @success,
        transaction_id: @transaction_id
      }.to_json
    end
  end
end
