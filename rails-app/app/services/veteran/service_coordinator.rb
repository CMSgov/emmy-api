module Veteran
  class ServiceCoordinator
    def initialize
      @client = Veteran::VaClient.new
    end

    def lookup_disability_rating(req_params)
      light = Stoplight("va-lookup-disability-rating")

      # In Stoplight 5.x, with_allowed_errors might not be on the light object directly
      # but it's part of the configuration.
      # The previous implementation used with_allowed_errors if it responded to it.
      if light.respond_to?(:with_allowed_errors)
        light = light.with_allowed_errors([ Veteran::NotFoundError ])
      end

      light.run do
        @client.lookup_disability_rating(req_params)
      end
    end
  end
end
