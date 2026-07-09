require_relative "../../application/port/payment_repository"
require_relative "../../domain/payment_reservation"

module PaymentService
  module Infrastructure
    module Persistence
      class SequelPaymentRepository < PaymentService::Application::Port::PaymentRepository
        def initialize(db)
          @db = db
        end

        # Mesma garantia do order-service: reserva de pagamento e evento de
        # outbox gravados na mesma transação. Se o commit falhar, nenhum dos
        # dois fica persistido.
        def save_with_outbox_event(reservation, event)
          @db.transaction do
            @db[:payments].insert(
              id: reservation.id,
              order_id: reservation.order_id,
              customer_id: reservation.customer_id,
              amount: reservation.amount,
              status: reservation.status,
              created_at: reservation.created_at
            )

            @db[:outbox_events].insert(
              id: event[:id],
              aggregate_id: event[:aggregate_id],
              event_type: event[:event_type],
              payload: Sequel.pg_jsonb(event[:payload]),
              published: false,
              created_at: reservation.created_at
            )
          end
        end

        def find_by_order_id(order_id)
          row = @db[:payments].where(order_id: order_id).first
          return nil unless row

          PaymentService::Domain::PaymentReservation.new(
            id: row[:id],
            order_id: row[:order_id],
            customer_id: row[:customer_id],
            amount: row[:amount].to_f,
            status: row[:status],
            created_at: row[:created_at]
          )
        end

        # Estorna a reserva (UPDATE, não INSERT) e grava o outbox event na mesma
        # transação — mesma garantia de atomicidade do save_with_outbox_event.
        def refund_with_outbox_event(reservation, event)
          @db.transaction do
            @db[:payments].where(order_id: reservation.order_id).update(status: reservation.status)

            @db[:outbox_events].insert(
              id: event[:id],
              aggregate_id: event[:aggregate_id],
              event_type: event[:event_type],
              payload: Sequel.pg_jsonb(event[:payload]),
              published: false,
              created_at: Time.now.utc
            )
          end
        end
      end
    end
  end
end
