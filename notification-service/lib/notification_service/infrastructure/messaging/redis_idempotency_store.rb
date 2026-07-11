require "redis"
require_relative "../../application/port/idempotency_store"

module NotificationService
  module Infrastructure
    module Messaging
      class RedisIdempotencyStore < NotificationService::Application::Port::IdempotencyStore
        TTL_SECONDS = 86_400 # 24h: tempo suficiente para cobrir reentregas do Kafka

        def initialize(redis)
          @redis = redis
        end

        # Reprocessamento de mensagem (retry do consumer, rebalance de partição)
        # não pode confirmar o mesmo pedido duas vezes nem publicar order.confirmed
        # repetido — é essa checagem que garante isso.
        def already_processed?(key)
          @redis.exists?("notification:processed:#{key}")
        end

        def mark_processed(key)
          @redis.set("notification:processed:#{key}", "1", ex: TTL_SECONDS)
        end
      end
    end
  end
end
