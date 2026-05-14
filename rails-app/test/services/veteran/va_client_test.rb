require 'test_helper'
require 'minitest/mock'
require 'net/http'
require 'stoplight'

module Veteran
  class VaClientTest < ActiveSupport::TestCase
    setup do
      @env_overrides = {
        'VA_BASE_URL' => 'https://example.test/va',
        'VA_TOKEN_URL' => 'https://example.test/token',
        'VA_CLIENT_ID' => 'client-id',
        'VA_AUD' => 'https://example.test/aud',
        'VA_PRIVATE_KEY_PATH' => 'test/fixtures/files/test.key',
        'SERVICE_VERSION' => '1.3.0',
        'ENVIRONMENT' => 'test'
      }
      @original_env = @env_overrides.keys.each_with_object({}) { |k, h| h[k] = ENV[k] }
      @env_overrides.each { |k, v| ENV[k] = v }

      # Create a dummy private key for testing
      @key = OpenSSL::PKey::RSA.new(2048)
      FileUtils.mkdir_p('test/fixtures/files')
      File.write('test/fixtures/files/test.key', @key.to_pem)

      @client = VaClient.new
      @coordinator = ServiceCoordinator.new
      @req_params = {
        firstName: 'Lynette',
        lastName: 'Oyola',
        dateOfBirth: '1988-10-24',
        ssn: '123456789'
      }
      Current.request_id = 'test-request-id'
    end

    teardown do
      File.delete('test/fixtures/files/test.key') if File.exist?('test/fixtures/files/test.key')
      @original_env.each { |k, v| ENV[k] = v }
    end

    private

    def stub_va_requests(oauth_resp, va_total_disability_resp = nil, rating_path: nil, total_disability_path: nil)
      http_mock_oauth = Minitest::Mock.new
      http_mock_oauth.expect :use_ssl=, true, [true]
      http_mock_oauth.expect :request, oauth_resp, [Net::HTTP::Post]

      http_mock_va = Minitest::Mock.new
      http_mock_va.expect :use_ssl=, true, [true]

      if va_total_disability_resp
        total_disability_matcher = ->(req) {
          req.is_a?(Net::HTTP::Post) && (total_disability_path.nil? || req.path == total_disability_path)
        }
        http_mock_va.expect :request, va_total_disability_resp do |req|
          total_disability_matcher.call(req)
        end
      end

      calls = 0
      Net::HTTP.stub :new, proc { |host, port|
        calls += 1
        calls == 1 ? http_mock_oauth : http_mock_va
      } do
        yield
      end

      http_mock_oauth.verify
      http_mock_va.verify
    end

    public

    test 'lookup_disability_rating success uses restricted endpoint when SSN is present' do
      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      va_rating_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      va_rating_response.instance_variable_set(:@read, true)
      va_rating_response.instance_variable_set(:@body, {
        data: {
          attributes: {
            combined_disability_rating: 100,
            combined_effective_date: '2023-01-01',
            legal_effective_date: '2023-01-01',
            individual_ratings: [
              { rating_end_date: '2024-01-01' },
              { rating_end_date: '2023-12-01' }
            ]
          }
        }
      }.to_json)

      va_total_disability_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      va_total_disability_response.instance_variable_set(:@read, true)
      va_total_disability_response.instance_variable_set(:@body, {
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
                       rating_path: '/va/restricted/disability_rating',
                       total_disability_path: '/va/restricted/permanent_and_total_disability') do
        response = @client.lookup_disability_rating(@req_params)
        assert_equal true, response.total_disability_status
        assert_equal true, response.total_disability_status_effective_date.present?
        assert_equal 'test-request-id', response.metadata[:transactionId]
      end
    end

    test 'lookup_disability_rating success uses standard endpoint when SSN is absent' do
      params_without_ssn = @req_params.except(:ssn)

      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      va_rating_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      va_rating_response.instance_variable_set(:@read, true)
      va_rating_response.instance_variable_set(:@body, {
        data: {
          attributes: {
            combined_disability_rating: 70,
            combined_effective_date: '2023-01-01',
            legal_effective_date: '2023-01-01'
          }
        }
      }.to_json)

      va_total_disability_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      va_total_disability_response.instance_variable_set(:@read, true)
      va_total_disability_response.instance_variable_set(:@body, {
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
                       rating_path: '/va/disability_rating',
                       total_disability_path: '/va/permanent_and_total_disability') do
        response = @client.lookup_disability_rating(params_without_ssn)
        assert_equal true, response.total_disability_status
        assert_equal true, response.total_disability_status_effective_date.present?
      end
    end

    test 'lookup_disability_rating raises NotFoundError on 404' do
      oauth_response = Net::HTTPSuccess.new('1.1', '200', 'OK')
      oauth_response.instance_variable_set(:@read, true)
      oauth_response.instance_variable_set(:@body, { access_token: 'fake-token' }.to_json)

      va_response = Net::HTTPNotFound.new('1.1', '404', 'Not Found')

      stub_va_requests(oauth_response, va_response) do
        assert_raises(Veteran::NotFoundError) do
          @client.lookup_disability_rating(@req_params)
        end
      end
    end
  end
end
