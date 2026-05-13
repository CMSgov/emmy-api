require 'test_helper'
require 'minitest/mock'
require 'net/http'
require 'stoplight'

module Education
  class NscClientTest < ActiveSupport::TestCase
    setup do
      @env_overrides = {
        'NSC_SUBMIT_URL' => 'https://example.test/submit',
        'NSC_ACCOUNT_ID' => '12345',
        'NSC_CLIENT_ID' => 'client-id',
        'NSC_CLIENT_SECRET' => 'client-secret',
        'NSC_TOKEN_URL' => 'https://example.test/token'
      }
      @original_env = @env_overrides.keys.each_with_object({}) { |k, h| h[k] = ENV[k] }
      @env_overrides.each { |k, v| ENV[k] = v }

      @client = NscClient.new
      @coordinator = ServiceCoordinator.new
      @req_body = {
        firstName: 'Lynette',
        lastName: 'Oyola',
        dateOfBirth: '1988-10-24'
      }
    end

    teardown do
      @original_env.each { |k, v| ENV[k] = v }
    end

    private

    def stub_nsc_requests(oauth_resp, submit_resp)
      http_mock_oauth = Minitest::Mock.new
      http_mock_oauth.expect :use_ssl=, true, [true]
      http_mock_oauth.expect :request, oauth_resp, [Net::HTTP::Post]

      http_mock_submit = Minitest::Mock.new
      http_mock_submit.expect :use_ssl=, true, [true]
      http_mock_submit.expect :request, submit_resp, [Net::HTTP::Post]

      # We need to stub Net::HTTP.new to return different mocks for different calls.
      # A simple way with Minitest::Mock is to expect :new and return mocks in order.
      # But Net::HTTP.stub :new stubs the class method.

      calls = 0
      Net::HTTP.stub :new, proc { |host, port|
        calls += 1
        calls == 1 ? http_mock_oauth : http_mock_submit
      } do
        yield
      end

      http_mock_oauth.verify
      http_mock_submit.verify
    end

    public

    test 'lookup_enrollment_status success with positive hit' do
      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      submit_response.instance_variable_set(:@read, true)
      submit_response.instance_variable_set(:@body, {
        status: { code: '0', message: 'Successful', severity: 'Info' },
        transactionDetails: { nscHit: 'Y', transactionStatus: 'CNF' },
        enrollmentDetails: [{ currentEnrollmentStatus: 'CC' }]
      }.to_json)

      stub_nsc_requests(oauth_response, submit_response) do
        result = @client.lookup_enrollment_status(@req_body).as_json

        assert_equal 'ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING', result[:enrollmentStatus]
        assert_equal 'NSC', result[:dataSource]
        assert_not_nil result[:metadata]
      end
    end

    test 'lookup_enrollment_status maps specific enrollment status' do
      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      submit_response.instance_variable_set(:@read, true)
      submit_response.instance_variable_set(:@body, {
        transactionDetails: { nscHit: 'Y' },
        enrollmentDetails: [{
          currentEnrollmentStatus: 'CC',
          officialSchoolName: 'University A',
          enrollmentData: [{ enrollmentStatus: 'H', termBeginDate: '2023-01-01', termEndDate: '2023-05-01' }]
        }]
      }.to_json)

      stub_nsc_requests(oauth_response, submit_response) do
        result = @client.lookup_enrollment_status(@req_body).as_json

        assert_equal 'HALF_TIME', result[:enrollmentStatus]
        assert_equal 1, result[:enrollmentDetails].size
        assert_equal 'HALF_TIME', result[:enrollmentDetails][0][:enrollmentStatus]
        assert_equal 'University A', result[:enrollmentDetails][0][:schoolName]
      end
    end

    test 'lookup_enrollment_status raises Not Found for no hit' do
      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      submit_response.instance_variable_set(:@read, true)
      submit_response.instance_variable_set(:@body, {
        transactionDetails: { nscHit: 'N' },
        enrollmentDetails: [{ currentEnrollmentStatus: 'CN' }]
      }.to_json)

      stub_nsc_requests(oauth_response, submit_response) do
        assert_raises(Education::NotFoundError) do
          @client.lookup_enrollment_status(@req_body)
        end
      end
    end

    test 'lookup_enrollment_status raises Not Found for currently not enrolled (CN)' do
      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      submit_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      submit_response.instance_variable_set(:@read, true)
      submit_response.instance_variable_set(:@body, {
        transactionDetails: { nscHit: 'Y' },
        enrollmentDetails: [{ currentEnrollmentStatus: 'CN' }]
      }.to_json)

      stub_nsc_requests(oauth_response, submit_response) do
        assert_raises(Education::NotFoundError) do
          @client.lookup_enrollment_status(@req_body)
        end
      end
    end

    test 'lookup_enrollment_status raises error for non-2xx response' do
      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      submit_response = Net::HTTPBadGateway.new('1.1', '502', 'Bad Gateway')

      stub_nsc_requests(oauth_response, submit_response) do
        assert_raises(StandardError, /NSC submit failed: status=502/) do
          @client.lookup_enrollment_status(@req_body)
        end
      end
    end
  end
end
