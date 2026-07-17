require "time"
require "securerandom"

module PaymentService
  module Domain
    class PaymentReservedEvent
      def self.from(reservation, saga_id:, correlation_id:)
        {
          eventType: "payment.reserved",
          version: 1,
          sagaId: saga_id,
          correlationId: correlation_id,
          occurredAt: reservation.created_at.iso8601,
          payload: {
            paymentId: reservation.id,
            orderId: reservation.order_id,
            customerId: reservation.customer_id,
            amount: reservation.amount
          }
        }
      end
    end
  end
end
