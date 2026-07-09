require "racecar"
require "json"
require_relative "../persistence/database"
require_relative "../persistence/sequel_payment_repository"
require_relative "../../application/use_case/refund_payment"

module PaymentService
  module Infrastructure
    module Messaging
      class InventoryFailedConsumer < Racecar::Consumer
        subscribes_to "inventory.failed"

        def initialize
          db = PaymentService::Infrastructure::Persistence.connect
          repository = PaymentService::Infrastructure::Persistence::SequelPaymentRepository.new(db)
          @use_case = PaymentService::Application::UseCase::RefundPayment.new(repository: repository)
        end

        def process(message)
          event = JSON.parse(message.value)

          @use_case.execute(
            inventory_failed_payload: event["payload"],
            saga_id: event["sagaId"],
            correlation_id: event["correlationId"]
          )

          puts "[inventory.failed] refunded payment for order #{event['payload']['orderId']}"
        rescue => e
          puts "[inventory.failed] failed to process message: #{e.message}"
          raise
        end
      end
    end
  end
end
