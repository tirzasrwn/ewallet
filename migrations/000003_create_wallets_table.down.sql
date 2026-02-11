-- Drop wallets indexes and table
DROP INDEX IF EXISTS idx_wallets_created_at;
DROP INDEX IF EXISTS idx_wallets_deleted_at;
DROP INDEX IF EXISTS idx_wallets_user_id;
DROP TABLE IF EXISTS wallets;
