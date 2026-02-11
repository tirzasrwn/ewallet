-- Rollback Demo Data
-- This migration removes all demo/test data

BEGIN;

-- Delete demo transactions
DELETE FROM transactions 
WHERE receiver_id IN (
    SELECT id FROM users WHERE email IN (
        'alice@example.com',
        'bob@example.com',
        'charlie@example.com',
        'diana@example.com',
        'eve@example.com'
    )
) OR sender_id IN (
    SELECT id FROM users WHERE email IN (
        'alice@example.com',
        'bob@example.com',
        'charlie@example.com',
        'diana@example.com',
        'eve@example.com'
    )
);

-- Delete demo wallets
DELETE FROM wallets 
WHERE user_id IN (
    SELECT id FROM users WHERE email IN (
        'alice@example.com',
        'bob@example.com',
        'charlie@example.com',
        'diana@example.com',
        'eve@example.com'
    )
);

-- Delete demo users
DELETE FROM users 
WHERE email IN (
    'alice@example.com',
    'bob@example.com',
    'charlie@example.com',
    'diana@example.com',
    'eve@example.com'
);

COMMIT;
