require "securerandom"
require_relative "../../domain/payment_reservation"
require_relative "../../domain/payment_reserved_event"
require_relative "../../domain/payment_failed_event"

module PaymentService
  module Application
    module UseCase
      class ReservePayment
        def initialize(repository:)
          @repository = repository
        end

        # order_created_payload é o payload já parseado do evento order.created;
        # saga_id/correlation_id vêm do envelope do evento, não são gerados aqui —
        # é isso que mantém o rastreamento coerente do início ao fim da saga.
        def execute(order_created_payload:, saga_id:, correlation_id:)
          reservation = build_reservation(order_created_payload)
          event = build_outbox_event(reservation, saga_id, correlation_id)

          @repository.save_with_outbox_event(reservation, event)
          reservation
        end

        private

        def build_reservation(payload)
          amount = payload["totalAmount"].to_f

          if amount.positive?
            PaymentService::Domain::PaymentReservation.reserve(
              id: SecureRandom.uuid,
              order_id: payload["orderId"],
              customer_id: payload["customerId"],
              amount: amount
            )
          else
            PaymentService::Domain::PaymentReservation.failed(
              id: SecureRandom.uuid,
              order_id: payload["orderId"],
              customer_id: payload["customerId"],
              amount: amount
            )
          end
        end

        def build_outbox_event(reservation, saga_id, correlation_id)
          domain_event = if reservation.reserved?
            PaymentService::Domain::PaymentReservedEvent.from(
              reservation, saga_id: saga_id, correlation_id: correlation_id
            )
          else
            PaymentService::Domain::PaymentFailedEvent.from(
              reservation, saga_id: saga_id, correlation_id: correlation_id, reason: "invalid order amount"
            )
          end

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
