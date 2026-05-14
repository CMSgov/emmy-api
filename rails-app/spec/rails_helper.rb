# This file is copied to spec/ when you run 'rails generate rspec:install'
require 'spec_helper'
ENV['RAILS_ENV'] ||= 'test'

# Mock AWS credentials and region to avoid MissingCredentialsError/MissingRegionError during tests
ENV['AWS_ACCESS_KEY_ID'] ||= 'testing'
ENV['AWS_SECRET_ACCESS_KEY'] ||= 'testing'
ENV['AWS_REGION'] ||= 'us-east-1'
ENV['DB_IAM_AUTH'] ||= 'false'

# CRITICAL: Stub PG::AWS_RDS_IAM as early as possible before Rails boots
module PG
  module AWS_RDS_IAM
    class AuthTokenGeneratorRegistry
      def self.instance; @instance ||= new; end
      def fetch(name); ->(*) { "dummy_token" }; end
      def reset; end
    end
    class AuthTokenInjector
      def inject_into_connection_string(conninfo); conninfo; end
    end
  end
end

require_relative '../config/environment'

# Monkeypatch ActiveRecord connection configuration to remove AWS IAM options in tests
module ActiveRecord
  module ConnectionAdapters
    class PostgreSQLAdapter
      class << self
        alias_method :original_new_client, :new_client
        def new_client(conninfo)
          if conninfo.is_a?(Hash)
            conninfo.delete(:aws_rds_iam_auth_token_generator)
            conninfo.delete(:aws_rds_iam_region)
          end
          original_new_client(conninfo)
        end
      end
    end
  end
end
# Prevent database truncation if the environment is production
abort("The Rails environment is running in production mode!") if Rails.env.production?

# FORCE pg-aws_rds_iam to be a no-op by monkeypatching the core generator
begin
  require 'pg/aws_rds_iam'
  module PG
    module AWS_RDS_IAM
      class AuthTokenGenerator
        def call(*)
          "dummy_token"
        end
      end
      class AuthTokenInjector
        def inject_into_connection_string(conninfo)
          conninfo
        end
      end
    end
  end
rescue LoadError
end

# Stub Aws::RDS::Client and other potential error sources
begin
  require 'aws-sdk-rds'
  module Aws
    module RDS
      class Client
        def self.new(*args, **kwargs)
          obj = allocate
          obj.instance_variable_set(:@config, Struct.new(:region).new('us-east-1'))
          obj
        end
        def initialize(*args, **kwargs); end
        def config; @config; end
      end
      class AuthTokenGenerator
        def initialize(*args, **kwargs); end
        def generate_auth_token(*args, **kwargs)
          "mocked_token"
        end
      end
    end
  end
rescue LoadError
end

begin
  require 'aws-sigv4'
  module Aws
    module Sigv4
      class Signer
        def initialize(*args, **kwargs); end
        def sign_request(*args, **kwargs)
          Struct.new(:headers).new({})
        end
      end
    end
  end
rescue LoadError
end

# Patch RegionalEndpoint to bypass MissingRegionError
begin
  require 'aws-sdk-core/plugins/regional_endpoint'
  module Aws
    module Plugins
      class RegionalEndpoint
        def after_initialize(client)
          # bypass
        end
      end
    end
  end
rescue LoadError
end

# Uncomment the line below in case you have `--require rails_helper` in the `.rspec` file
# that will avoid rails generators crashing because migrations haven't been run yet
# return unless Rails.env.test?
require 'rspec/rails'
# Add additional requires below this line. Rails is not loaded until this point!

# Requires supporting ruby files with custom matchers and macros, etc, in
# spec/support/ and its subdirectories. Files matching `spec/**/*_spec.rb` are
# run as spec files by default. This means that files in spec/support that end
# in _spec.rb will both be required and run as specs, causing the specs to be
# run twice. It is recommended that you do not name files matching this glob to
# end with _spec.rb. You can configure this pattern with the --pattern
# option on the command line or in ~/.rspec, .rspec or `.rspec-local`.
#
# The following line is provided for convenience purposes. It has the downside
# of increasing the boot-up time by auto-requiring all files in the support
# directory. Alternatively, in the individual `*_spec.rb` files, manually
# require only the support files necessary.
#
# Rails.root.glob('spec/support/**/*.rb').sort_by(&:to_s).each { |f| require f }
Dir[Rails.root.join('spec/support/**/*.rb')].sort.each { |f| require f }

# Checks for pending migrations and applies them before tests are run.
# If you are not using ActiveRecord, you can remove these lines.
# begin
#   ActiveRecord::Migration.maintain_test_schema!
# rescue ActiveRecord::PendingMigrationError => e
#   abort e.to_s.strip
# end
RSpec.configure do |config|
  # Remove this line if you're not using ActiveRecord or ActiveRecord fixtures
  config.fixture_paths = [
    Rails.root.join('spec/fixtures')
  ]

  # If you're not using ActiveRecord, or you'd prefer not to run each of your
  # examples within a transaction, remove the following line or assign false
  # instead of true.
  config.use_transactional_fixtures = true

  # You can uncomment this line to turn off ActiveRecord support entirely.
  # config.use_active_record = false

  # RSpec Rails uses metadata to mix in different behaviours to your tests,
  # for example enabling you to call `get` and `post` in request specs. e.g.:
  #
  #     RSpec.describe UsersController, type: :request do
  #       # ...
  #     end
  #
  # The different available types are documented in the features, such as in
  # https://rspec.info/features/7-1/rspec-rails
  #
  # You can also this infer these behaviours automatically by location, e.g.
  # /spec/models would pull in the same behaviour as `type: :model` but this
  # behaviour is considered legacy and will be removed in a future version.
  #
  # To enable this behaviour uncomment the line below.
  # config.infer_spec_type_from_file_location!

  # Filter lines from Rails gems in backtraces.
  config.filter_rails_from_backtrace!
  # arbitrary gems may also be filtered via:
  # config.filter_gems_from_backtrace("gem name")
end
