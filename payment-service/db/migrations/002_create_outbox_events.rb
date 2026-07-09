Sequel.migration do
  change do
    create_table(:outbox_events) do
      column :id, "uuid", primary_key: true
      column :aggregate_id, "uuid", null: false
      String :event_type, null: false
      column :payload, "jsonb", null: false
      TrueClass :published, default: false, null: false
      DateTime :created_at, null: false

      index [:published, :created_at]
    end
  end
end
