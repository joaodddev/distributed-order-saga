Sequel.migration do
  change do
    create_table(:payments) do
      column :id, "uuid", primary_key: true
      column :order_id, "uuid", null: false
      column :customer_id, "uuid", null: false
      Numeric :amount, size: [12, 2], null: false
      String :status, null: false
      DateTime :created_at, null: false

      index :order_id
    end
  end
end
