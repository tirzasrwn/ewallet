-- Drop indexes
DROP INDEX IF EXISTS idx_wallets_deleted_at;
DROP INDEX IF EXISTS idx_wallets_user_id;

-- Drop wallets table
DROP TABLE IF EXISTS wallets;
