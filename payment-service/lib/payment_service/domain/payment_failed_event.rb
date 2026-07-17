require "time"

module PaymentService
  module Domain
    class PaymentFailedEvent
      def self.from(reservation, saga_id:, correlation_id:, reason:)
        {
          eventType: "payment.failed",
          version: 1,
          sagaId: saga_id,
          correlationId: correlation_id,
          occurredAt: reservation.created_at.iso8601,
          payload: {
            paymentId: reservation.id,
            orderId: reservation.order_id,
            customerId: reservation.customer_id,
            reason: reason
          }
        }
      end
    end
  end
end
