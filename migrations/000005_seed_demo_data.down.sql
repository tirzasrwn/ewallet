BEGIN;
-- Remove demo transactions
DELETE FROM transactions
WHERE sender_id IN (
  SELECT id FROM users WHERE email IN (
    'alice@example.com','bob@example.com','charlie@example.com','diana@example.com','eve@example.com'
  )
) OR receiver_id IN (
  SELECT id FROM users WHERE email IN (
    'alice@example.com','bob@example.com','charlie@example.com','diana@example.com','eve@example.com'
  )
);

-- Remove demo wallets
DELETE FROM wallets
WHERE user_id IN (
  SELECT id FROM users WHERE email IN (
    'alice@example.com','bob@example.com','charlie@example.com','diana@example.com','eve@example.com'
  )
);

-- Remove demo users
DELETE FROM users
WHERE email IN (
  'alice@example.com','bob@example.com','charlie@example.com','diana@example.com','eve@example.com'
);

COMMIT;
