CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_id UUID,
  receiver_id UUID NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  type VARCHAR(20) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES users(id),
  CONSTRAINT fk_receiver FOREIGN KEY (receiver_id) REFERENCES users(id)
);
