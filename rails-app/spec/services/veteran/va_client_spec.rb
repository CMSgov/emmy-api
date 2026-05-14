require 'rails_helper'
require 'net/http'

module Veteran
  RSpec.describe VaClient, type: :service do
    let(:env_overrides) do
      {
        'VA_BASE_URL' => 'https://example.test/va',
        'VA_TOKEN_URL' => 'https://example.test/token',
        'VA_CLIENT_ID' => 'client-id',
        'VA_AUD' => 'https://example.test/aud',
        'VA_PRIVATE_KEY_PATH' => 'spec/fixtures/files/test.key',
        'SERVICE_VERSION' => '1.3.0',
        'ENVIRONMENT' => 'test'
      }
    end

    around do |example|
      # Create a dummy private key for testing
      key = OpenSSL::PKey::RSA.new(2048)
      FileUtils.mkdir_p('spec/fixtures/files')
      File.write('spec/fixtures/files/test.key', key.to_pem)

      ClimateControl.modify(env_overrides) do
        example.run
      end

      File.delete('spec/fixtures/files/test.key') if File.exist?('spec/fixtures/files/test.key')
    end

    let(:client) { VaClient.new }
    let(:req_params) do
      {
        firstName: 'Lynette',
        lastName: 'Oyola',
        dateOfBirth: '1988-10-24',
        ssn: '123456789'
      }
    end

    before do
      Current.request_id = 'test-request-id'
    end

    def stub_va_requests(oauth_resp, va_total_disability_resp = nil, total_disability_path: nil)
      http_mock_oauth = instance_double(Net::HTTP)
      allow(http_mock_oauth).to receive(:use_ssl=).with(true)
      allow(http_mock_oauth).to receive(:request).with(instance_of(Net::HTTP::Post)).and_return(oauth_resp)

      http_mock_va = instance_double(Net::HTTP)
      allow(http_mock_va).to receive(:use_ssl=).with(true)

      if va_total_disability_resp
        allow(http_mock_va).to receive(:request) do |req|
          if total_disability_path.nil? || req.path == total_disability_path
            va_total_disability_resp
          else
            # Fallback for rating request if needed, though the client calls it in order
            nil
          end
        end
      end

      expect(Net::HTTP).to receive(:new).and_return(http_mock_oauth, http_mock_va)
    end

    describe '#lookup_disability_rating' do
      it 'success uses restricted endpoint when SSN is present' do
        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        va_total_disability_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(va_total_disability_response).to receive(:body).and_return({
          data: {
            total_disability: {
              status: true,
              effective_date: "2023-01-01"
            },
            permanent_and_total: {
              service_connected_status: false,
              pension_award_status: false
            }
          }
        }.to_json)

        stub_va_requests(oauth_response, va_total_disability_response,
                         total_disability_path: '/va/restricted/permanent_and_total_disability')

        response = client.lookup_disability_rating(req_params)
        expect(response.total_disability_status).to be true
        expect(response.total_disability_status_effective_date).to be_present
        expect(response.metadata[:transactionId]).to eq('test-request-id')
      end

      it 'success uses standard endpoint when SSN is absent' do
        params_without_ssn = req_params.except(:ssn)

        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        va_total_disability_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(va_total_disability_response).to receive(:body).and_return({
          data: {
            total_disability: {
              status: true,
              effective_date: "2023-01-01"
            },
            permanent_and_total: {
              service_connected_status: true,
              pension_award_status: true
            }
          }
        }.to_json)

        stub_va_requests(oauth_response, va_total_disability_response,
                         total_disability_path: '/va/permanent_and_total_disability')

        response = client.lookup_disability_rating(params_without_ssn)
        expect(response.total_disability_status).to be true
        expect(response.total_disability_status_effective_date).to be_present
      end

      it 'raises NotFoundError on 404' do
        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        va_response = Net::HTTPNotFound.new('1.1', '404', 'Not Found')

        stub_va_requests(oauth_response, va_response)

        expect {
          client.lookup_disability_rating(req_params)
        }.to raise_error(Veteran::NotFoundError)
      end
    end
  end
end
