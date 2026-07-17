require "securerandom"
require "time"

module NotificationService
  module Domain
    class OrderConfirmedEvent
      def self.from(order_id:, customer_id:, saga_id:, correlation_id:)
        {
          eventType: "order.confirmed",
          version: 1,
          sagaId: saga_id,
          correlationId: correlation_id,
          occurredAt: Time.now.utc.iso8601,
          payload: {
            orderId: order_id,
            customerId: customer_id
          }
        }
      end
    end
  end
end
