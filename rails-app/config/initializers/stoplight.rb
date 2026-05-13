require 'stoplight'

# Custom configuration for Stoplight can be added here
# By default Stoplight uses an in-memory data store.
if defined?(Stoplight)
  # Try to configure the light if the module is available
  begin
    # Check if the class is actually available
    if Stoplight.const_defined?(:Light)
      Stoplight::Light.default_error_handler = lambda do |error, handle|
        # We want to raise specific errors so they can be handled by the controller
        # Stoplight::Error::RedLight is raised when the circuit is open
        raise error
      end
    end
  rescue NameError
    # Fallback or log if needed
  end
end
