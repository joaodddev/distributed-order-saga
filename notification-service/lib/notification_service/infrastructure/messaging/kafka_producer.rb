require "rdkafka"

module NotificationService
  module Infrastructure
    module Messaging
      class KafkaProducer
        def initialize(brokers:)
          @producer = Rdkafka::Config.new("bootstrap.servers" => brokers).producer
        end

        def publish(topic:, key:, payload:)
          @producer.produce(topic: topic, payload: payload, key: key).wait
        end

        def close
          @producer.close
        end
      end
    end
  end
end
