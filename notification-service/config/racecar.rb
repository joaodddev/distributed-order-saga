require_relative "../lib/notification_service/infrastructure/messaging/inventory_reserved_consumer"

Racecar.configure do |config|
  config.brokers = [ENV.fetch("KAFKA_BROKERS", "localhost:9092")]
  config.client_id = "notification-service"
  config.group_id = "notification-service-group"
end
