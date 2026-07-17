require "time"

module PaymentService
  module Domain
    class PaymentRefundedEvent
      def self.from(reservation, saga_id:, correlation_id:, reason:)
        {
          eventType: "payment.refunded",
          version: 1,
          sagaId: saga_id,
          correlationId: correlation_id,
          occurredAt: Time.now.utc.iso8601,
          payload: {
            paymentId: reservation.id,
            orderId: reservation.order_id,
            customerId: reservation.customer_id,
            amount: reservation.amount,
            reason: reason
          }
        }
      end
    end
  end
end
