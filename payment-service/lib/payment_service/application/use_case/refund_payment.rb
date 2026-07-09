require "securerandom"
require_relative "../../domain/payment_refunded_event"

module PaymentService
  module Application
    module UseCase
      class RefundPayment
        def initialize(repository:)
          @repository = repository
        end

        # inventory_failed_payload vem do evento inventory.failed:
        # { orderId, customerId, reason }
        def execute(inventory_failed_payload:, saga_id:, correlation_id:)
          reservation = @repository.find_by_order_id(inventory_failed_payload["orderId"])
          return unless reservation # idempotência: se já não existe reserva, não há o que estornar

          reservation.refund!

          event = build_outbox_event(reservation, saga_id, correlation_id, inventory_failed_payload["reason"])
          @repository.refund_with_outbox_event(reservation, event)
        end

        private

        def build_outbox_event(reservation, saga_id, correlation_id, reason)
          domain_event = PaymentService::Domain::PaymentRefundedEvent.from(
            reservation, saga_id: saga_id, correlation_id: correlation_id, reason: reason
          )

          {
            id: SecureRandom.uuid,
            aggregate_id: reservation.id,
            event_type: domain_event[:eventType],
            payload: domain_event
          }
        end
      end
    end
  end
end
