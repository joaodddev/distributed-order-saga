require "spec_helper"
require "payment_service/application/use_case/reserve_payment"

RSpec.describe PaymentService::Application::UseCase::ReservePayment do
  # double em memória, mesma ideia do fake em Go
  class FakeRepository
    attr_reader :saved_reservation, :saved_event

    def save_with_outbox_event(reservation, event)
      @saved_reservation = reservation
      @saved_event = event
    end
  end

  it "reserves payment and builds payment.reserved event for a valid order" do
    repository = FakeRepository.new
    use_case = described_class.new(repository: repository)

    use_case.execute(
      order_created_payload: {
        "orderId" => SecureRandom.uuid,
        "customerId" => SecureRandom.uuid,
        "totalAmount" => 99.80
      },
      saga_id: "saga-1",
      correlation_id: "corr-1"
    )

    expect(repository.saved_reservation.status).to eq("RESERVED")
    expect(repository.saved_event[:event_type]).to eq("payment.reserved")
  end

  it "fails payment and builds payment.failed event for zero amount" do
    repository = FakeRepository.new
    use_case = described_class.new(repository: repository)

    use_case.execute(
      order_created_payload: {
        "orderId" => SecureRandom.uuid,
        "customerId" => SecureRandom.uuid,
        "totalAmount" => 0
      },
      saga_id: "saga-1",
      correlation_id: "corr-1"
    )

    expect(repository.saved_reservation.status).to eq("FAILED")
    expect(repository.saved_event[:event_type]).to eq("payment.failed")
  end
end
