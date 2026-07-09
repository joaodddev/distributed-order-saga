module PaymentService
  module Domain
    class PaymentReservation
      STATUSES = %w[RESERVED FAILED REFUNDED].freeze

      attr_reader :id, :order_id, :customer_id, :amount, :status, :created_at

      def initialize(id:, order_id:, customer_id:, amount:, status:, created_at: Time.now.utc)
        @id = id
        @order_id = order_id
        @customer_id = customer_id
        @amount = amount
        @status = status
        @created_at = created_at
      end

      def self.reserve(id:, order_id:, customer_id:, amount:)
        new(id: id, order_id: order_id, customer_id: customer_id, amount: amount, status: "RESERVED")
      end

      def self.failed(id:, order_id:, customer_id:, amount:)
        new(id: id, order_id: order_id, customer_id: customer_id, amount: amount, status: "FAILED")
      end

      def refund!
        @status = "REFUNDED"
      end

      def reserved?
        status == "RESERVED"
      end

      def failed?
        status == "FAILED"
      end
    end
  end
end
