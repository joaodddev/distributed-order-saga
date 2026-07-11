require "opentelemetry/sdk"
require "opentelemetry/exporter/otlp"

module NotificationService
  module Infrastructure
    module Observability
      def self.setup(service_name:, endpoint:)
        OpenTelemetry::SDK.configure do |c|
          c.service_name = service_name
          c.add_span_processor(
            OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(
              OpenTelemetry::Exporter::OTLP::Exporter.new(endpoint: endpoint)
            )
          )
        end
      end

      def self.tracer
        OpenTelemetry.tracer_provider.tracer("notification-service")
      end
    end
  end
end
