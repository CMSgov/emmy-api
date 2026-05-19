require "datadog"

Datadog.configure do |c|
  # Default tracing to disabled in development unless explicitly enabled via environment variable
  default_enabled = Rails.env.production? ? "true" : "false"
  c.tracing.enabled = ENV.fetch("DD_TRACE_ENABLED", default_enabled).to_s.downcase == "true"

  # This will trace HTTP requests, Database queries, Redis, etc. automatically
  c.tracing.instrument :rails, service_name: ENV.fetch("DD_SERVICE", "emmy-rails-app")
  c.tracing.instrument :http, service_name: ENV.fetch("DD_SERVICE", "emmy-rails-app")
  c.tracing.instrument :shoryuken if defined?(Shoryuken)

  # Enable profiling if needed (requires libdatadog in the image)
  c.profiling.enabled = ENV.fetch("DD_PROFILING_ENABLED", "false") == "true"
end
