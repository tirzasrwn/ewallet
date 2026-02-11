-- Drop transaction indexes and table
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_deleted_at;
DROP INDEX IF EXISTS idx_transactions_sender_id;
DROP INDEX IF EXISTS idx_transactions_receiver_id;
DROP TABLE IF EXISTS transactions;
