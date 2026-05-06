require 'test_helper'

module Reporting
  class ReportJobTest < ActiveSupport::TestCase
    test 'performs and creates an ApiEvent' do
      body = {
        timestamp: Time.now.iso8601,
        endpoint: '/api/v0/test',
        data_source: 'test-source',
        client_id: 'test-client',
        status_code: 200,
        success: true
      }.to_json

      assert_difference 'ApiEvent.count', 1 do
        Reporting::ReportJob.new.perform(nil, body)
      end

      event = ApiEvent.last
      assert_equal '/api/v0/test', event.endpoint
      assert_equal 'test-source', event.data_source
      assert_equal 'test-client', event.client_id
      assert_equal 200, event.status_code
      assert_equal true, event.success
    end

    test 'logs error and raises if JSON is invalid' do
      assert_raises(JSON::ParserError) do
        Reporting::ReportJob.new.perform(nil, 'invalid json')
      end
    end
  end
end
