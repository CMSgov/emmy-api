require 'rails_helper'

module Reporting
  RSpec.describe ReportJob, type: :job do
    describe '#perform' do
      it 'performs and creates an ApiEvent' do
        body = {
          timestamp: Time.now.iso8601,
          endpoint: '/api/v0/test',
          data_source: 'test-source',
          client_id: 'test-client',
          status_code: 200,
          success: true
        }.to_json

        expect {
          ReportJob.new.perform(nil, body)
        }.to change(ApiEvent, :count).by(1)

        event = ApiEvent.last
        expect(event.endpoint).to eq('/api/v0/test')
        expect(event.data_source).to eq('test-source')
        expect(event.client_id).to eq('test-client')
        expect(event.status_code).to eq(200)
        expect(event.success).to be true
      end

      it 'logs error and raises if JSON is invalid' do
        expect {
          ReportJob.new.perform(nil, 'invalid json')
        }.to raise_error(JSON::ParserError)
      end
    end
  end
end
