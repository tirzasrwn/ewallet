-- Create wallets table with UUIDs
CREATE TABLE IF NOT EXISTS wallets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  balance DECIMAL(15,2) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT fk_wallet_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes for query performance
CREATE UNIQUE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_created_at ON wallets(created_at);
CREATE INDEX idx_wallets_deleted_at ON wallets(deleted_at);
