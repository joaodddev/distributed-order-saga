require "json"

module PaymentService
  module Outbox
    class Relay
      def initialize(db:, producer:, interval: 2, batch_size: 20)
        @db = db
        @producer = producer
        @interval = interval
        @batch_size = batch_size
      end

      def start
        loop do
          process_batch
          sleep @interval
        end
      end

      private

      def process_batch
        events = @db[:outbox_events]
                   .where(published: false)
                   .order(:created_at)
                   .limit(@batch_size)
                   .all

        events.each do |event|
          @producer.publish(
            topic: event[:event_type],
            key: event[:aggregate_id],
            payload: event[:payload].to_json
          )

          @db[:outbox_events].where(id: event[:id]).update(published: true)
        rescue => e
          puts "outbox relay: failed to publish event #{event[:id]}: #{e.message}"
        end
      end
    end
  end
end
