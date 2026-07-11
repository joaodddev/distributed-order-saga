require "securerandom"
require_relative "../../domain/order_confirmed_event"

module NotificationService
  module Application
    module UseCase
      class ConfirmOrder
        def initialize(idempotency_store:, publisher:)
          @idempotency_store = idempotency_store
          @publisher = publisher
        end

        # inventory_reserved_payload vem do evento inventory.reserved:
        # { reservationId, orderId, customerId }
        def execute(inventory_reserved_payload:, saga_id:, correlation_id:)
          order_id = inventory_reserved_payload["orderId"]

          return if @idempotency_store.already_processed?(order_id)

          event = NotificationService::Domain::OrderConfirmedEvent.from(
            order_id: order_id,
            customer_id: inventory_reserved_payload["customerId"],
            saga_id: saga_id,
            correlation_id: correlation_id
          )

          @publisher.publish(
            topic: "order.confirmed",
            key: order_id,
            payload: event.to_json
          )

          @idempotency_store.mark_processed(order_id)
        end
      end
    end
  end
end
