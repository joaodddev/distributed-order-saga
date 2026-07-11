require "racecar"
require "json"
require_relative "../persistence/database"
require_relative "../persistence/sequel_payment_repository"
require_relative "../../application/use_case/reserve_payment"
require_relative "../observability/tracer"

module PaymentService
  module Infrastructure
    module Messaging
      class OrderCreatedConsumer < Racecar::Consumer
        subscribes_to "order.created"

        def initialize
          PaymentService::Infrastructure::Observability.setup(
            service_name: "payment-service",
            endpoint: ENV.fetch("OTEL_COLLECTOR_ENDPOINT", "http://localhost:4318/v1/traces")
          )

          db = PaymentService::Infrastructure::Persistence.connect
          repository = PaymentService::Infrastructure::Persistence::SequelPaymentRepository.new(db)
          @use_case = PaymentService::Application::UseCase::ReservePayment.new(repository: repository)
        end

        def process(message)
          event = JSON.parse(message.value)

          @use_case.execute(
            order_created_payload: event["payload"],
            saga_id: event["sagaId"],
            correlation_id: event["correlationId"]
          )

          puts "[order.created] processed order #{event['payload']['orderId']}"
        rescue => e
          # Deixa a exceção subir: Racecar não commita o offset e reentrega
          # a mensagem. Prefiro reprocessar (idempotência fica pra fase de
          # compensação) do que perder o evento silenciosamente.
          puts "[order.created] failed to process message: #{e.message}"
          raise
        end
      end
    end
  end
end
