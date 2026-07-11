require "securerandom"
require_relative "../../domain/order_confirmed_event"
require_relative "../../infrastructure/observability/tracer"

module NotificationService
  module Application
    module UseCase
      class ConfirmOrder
        def initialize(idempotency_store:, publisher:)
          @idempotency_store = idempotency_store
          @publisher = publisher
        end

        def execute(inventory_reserved_payload:, saga_id:, correlation_id:)
          NotificationService::Infrastructure::Observability.tracer.in_span("ConfirmOrder.execute") do |span|
            order_id = inventory_reserved_payload["orderId"]
            span.set_attribute("saga.id", saga_id)
            span.set_attribute("order.id", order_id)

            if @idempotency_store.already_processed?(order_id)
              span.set_attribute("order.already_confirmed", true)
              next
            end

            event = NotificationService::Domain::OrderConfirmedEvent.from(
              order_id: order_id,
              customer_id: inventory_reserved_payload["customerId"],
              saga_id: saga_id,
              correlation_id: correlation_id
            )

            @publisher.publish(topic: "order.confirmed", key: order_id, payload: event.to_json)
            @idempotency_store.mark_processed(order_id)

            span.set_attribute("order.confirmed", true)
          end
        end
      end
    end
  end
end
