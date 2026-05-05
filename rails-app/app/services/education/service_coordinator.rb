module Education
  class ServiceCoordinator
    def initialize
      @client = Education::NscClient.new
    end

    def lookup_enrollment_status(req_params)
      light = Stoplight("nsc-lookup-enrollment-status")

      # We don't want to trip the circuit for NotFoundError
      if light.respond_to?(:with_allowed_errors)
        light = light.with_allowed_errors([Education::NotFoundError])
      end

      light.run do
        @client.lookup_enrollment_status(req_params)
      end
    end
  end
end
