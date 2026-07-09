module PaymentService
  module Domain
    class InvalidAmountError < StandardError; end

    class PaymentReservation
      STATUSES = %w[RESERVED FAILED REFUNDED].freeze

      attr_reader :id, :order_id, :customer_id, :amount, :status, :created_at

      def initialize(id:, order_id:, customer_id:, amount:, status: "RESERVED", created_at: Time.now.utc)
        raise InvalidAmountError, "amount must be greater than zero" unless amount.positive?

        @id = id
        @order_id = order_id
        @customer_id = customer_id
        @amount = amount
        @status = status
        @created_at = created_at
      end

      def fail!
        @status = "FAILED"
      end

      def refund!
        @status = "REFUNDED"
      end

      def reserved?
        status == "RESERVED"
      end
    end
  end
end
