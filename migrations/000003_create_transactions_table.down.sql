-- Drop indexes
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_deleted_at;
DROP INDEX IF EXISTS idx_transactions_receiver_id;
DROP INDEX IF EXISTS idx_transactions_sender_id;

-- Drop transactions table
DROP TABLE IF EXISTS transactions;
