require 'rails_helper'
require 'net/http'

module Education
  RSpec.describe NscClient, type: :service do
    let(:env_overrides) do
      {
        'NSC_SUBMIT_URL' => 'https://example.test/submit',
        'NSC_ACCOUNT_ID' => '12345',
        'NSC_CLIENT_ID' => 'client-id',
        'NSC_CLIENT_SECRET' => 'client-secret',
        'NSC_TOKEN_URL' => 'https://example.test/token'
      }
    end

    around do |example|
      ClimateControl.modify(env_overrides) do
        example.run
      end
    end

    let(:client) { NscClient.new }
    let(:req_body) do
      {
        firstName: 'Lynette',
        lastName: 'Oyola',
        dateOfBirth: '1988-10-24'
      }
    end

    def stub_nsc_requests(oauth_resp, submit_resp)
      http_mock_oauth = instance_double(Net::HTTP)
      allow(http_mock_oauth).to receive(:use_ssl=).with(true)
      allow(http_mock_oauth).to receive(:request).with(instance_of(Net::HTTP::Post)).and_return(oauth_resp)

      http_mock_submit = instance_double(Net::HTTP)
      allow(http_mock_submit).to receive(:use_ssl=).with(true)
      allow(http_mock_submit).to receive(:request).with(instance_of(Net::HTTP::Post)).and_return(submit_resp)

      expect(Net::HTTP).to receive(:new).and_return(http_mock_oauth, http_mock_submit)
    end

    describe '#lookup_enrollment_status' do
      it 'success with positive hit' do
        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(submit_response).to receive(:body).and_return({
          status: { code: '0', message: 'Successful', severity: 'Info' },
          transactionDetails: { nscHit: 'Y', transactionStatus: 'CNF' },
          enrollmentDetails: [ { currentEnrollmentStatus: 'CC' } ]
        }.to_json)

        stub_nsc_requests(oauth_response, submit_response)

        result = client.lookup_enrollment_status(req_body).as_json

        expect(result[:enrollmentStatus]).to eq('ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING')
        expect(result[:dataSource]).to eq('NSC')
        expect(result[:metadata]).not_to be_nil
      end

      it 'maps specific enrollment status' do
        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(submit_response).to receive(:body).and_return({
          transactionDetails: { nscHit: 'Y' },
          enrollmentDetails: [ {
            currentEnrollmentStatus: 'CC',
            officialSchoolName: 'University A',
            schoolCode: 117100,
            enrollmentData: [ { enrollmentStatus: 'H', termBeginDate: '2023-01-01', termEndDate: '2023-05-01' } ]
          } ]
        }.to_json)

        stub_nsc_requests(oauth_response, submit_response)

        result = client.lookup_enrollment_status(req_body).as_json

        expect(result[:enrollmentStatus]).to eq('HALF_TIME')
        expect(result[:enrollmentDetails].size).to eq(1)
        expect(result[:enrollmentDetails][0][:enrollmentStatus]).to eq('HALF_TIME')
        expect(result[:enrollmentDetails][0][:schoolName]).to eq('University A')
      end

      it 'raises Not Found for no hit' do
        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(submit_response).to receive(:body).and_return({
          transactionDetails: { nscHit: 'N' },
          enrollmentDetails: [ { currentEnrollmentStatus: 'CN' } ]
        }.to_json)

        stub_nsc_requests(oauth_response, submit_response)

        expect {
          client.lookup_enrollment_status(req_body)
        }.to raise_error(Education::NotFoundError)
      end

      it 'raises Not Found for currently not enrolled (CN)' do
        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(submit_response).to receive(:body).and_return({
          transactionDetails: { nscHit: 'Y' },
          enrollmentDetails: [ { currentEnrollmentStatus: 'CN' } ]
        }.to_json)

        stub_nsc_requests(oauth_response, submit_response)

        expect {
          client.lookup_enrollment_status(req_body)
        }.to raise_error(Education::NotFoundError)
      end

      it 'raises error for non-2xx response' do
        oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
        allow(oauth_response).to receive(:body).and_return({ access_token: 'fake-token' }.to_json)

        submit_response = Net::HTTPBadGateway.new('1.1', '502', 'Bad Gateway')
        submit_response.instance_variable_set(:@read, true)
        submit_response.body = '{"error":"upstream service unavailable"}'

        stub_nsc_requests(oauth_response, submit_response)

        expect {
          client.lookup_enrollment_status(req_body)
        }.to raise_error(StandardError, /NSC submit failed: status=502/)
      end
    end
  end
end
