module PaymentService
  module Application
    module Port
      # Porta de saída. Qualquer adapter (ex: SequelPaymentRepository) deve
      # implementar esses dois métodos com a mesma assinatura — Ruby não tem
      # interface formal, então isso aqui funciona como contrato + documentação.
      class PaymentRepository
        def save_with_outbox_event(reservation, event)
          raise NotImplementedError
        end

        def find_by_order_id(order_id)
          raise NotImplementedError
        end
      end
    end
  end
end
