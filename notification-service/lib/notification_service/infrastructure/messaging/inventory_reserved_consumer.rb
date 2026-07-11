require "racecar"
require "json"
require "redis"
require_relative "redis_idempotency_store"
require_relative "kafka_producer"
require_relative "../../application/use_case/confirm_order"

module NotificationService
  module Infrastructure
    module Messaging
      class InventoryReservedConsumer < Racecar::Consumer
        subscribes_to "inventory.reserved"

        def initialize
          redis = Redis.new(url: ENV.fetch("REDIS_URL", "redis://localhost:6379/0"))
          idempotency_store = RedisIdempotencyStore.new(redis)
          publisher = KafkaProducer.new(brokers: ENV.fetch("KAFKA_BROKERS", "localhost:9092"))

          @use_case = NotificationService::Application::UseCase::ConfirmOrder.new(
            idempotency_store: idempotency_store,
            publisher: publisher
          )
        end

        def process(message)
          event = JSON.parse(message.value)

          @use_case.execute(
            inventory_reserved_payload: event["payload"],
            saga_id: event["sagaId"],
            correlation_id: event["correlationId"]
          )

          puts "[inventory.reserved] confirmed order #{event['payload']['orderId']}"
        rescue => e
          puts "[inventory.reserved] failed to process message: #{e.message}"
          raise
        end
      end
    end
  end
end
