require "spec_helper"
require "notification_service/application/use_case/confirm_order"

RSpec.describe NotificationService::Application::UseCase::ConfirmOrder do
  class FakeIdempotencyStore
    def initialize(processed: false)
      @processed = processed
    end

    def already_processed?(key)
      @processed
    end

    def mark_processed(key); end
  end

  class FakePublisher
    attr_reader :published_topic, :published_payload

    def publish(topic:, key:, payload:)
      @published_topic = topic
      @published_payload = payload
    end
  end

  it "publishes order.confirmed when not previously processed" do
    publisher = FakePublisher.new
    use_case = described_class.new(
      idempotency_store: FakeIdempotencyStore.new(processed: false),
      publisher: publisher
    )

    use_case.execute(
      inventory_reserved_payload: { "orderId" => SecureRandom.uuid, "customerId" => SecureRandom.uuid },
      saga_id: "saga-1",
      correlation_id: "corr-1"
    )

    expect(publisher.published_topic).to eq("order.confirmed")
  end

  it "does not publish again when already processed (idempotency)" do
    publisher = FakePublisher.new
    use_case = described_class.new(
      idempotency_store: FakeIdempotencyStore.new(processed: true),
      publisher: publisher
    )

    use_case.execute(
      inventory_reserved_payload: { "orderId" => SecureRandom.uuid, "customerId" => SecureRandom.uuid },
      saga_id: "saga-1",
      correlation_id: "corr-1"
    )

    expect(publisher.published_topic).to be_nil
  end
end
