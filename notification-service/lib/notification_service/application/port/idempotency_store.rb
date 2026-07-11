module NotificationService
  module Application
    module Port
      class IdempotencyStore
        def already_processed?(key)
          raise NotImplementedError
        end

        def mark_processed(key)
          raise NotImplementedError
        end
      end
    end
  end
end
