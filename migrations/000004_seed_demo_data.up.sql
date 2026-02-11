-- Seed Demo Data
-- This migration adds demo/test data for development and testing

BEGIN;

-- Insert demo users
-- Password for all users: "password123" (bcrypt hashed)
-- Hash generated with: go run scripts/generate_hash.go password123
INSERT INTO users (name, email, password, created_at, updated_at) VALUES
('Alice Johnson', 'alice@example.com', '$2a$10$oRkxtds7C3lUdM2ZahxVmu5ZDTM.Cqj6e6Tc3gJzA.Vr/d/CWp8be', NOW(), NOW()),
('Bob Smith', 'bob@example.com', '$2a$10$oRkxtds7C3lUdM2ZahxVmu5ZDTM.Cqj6e6Tc3gJzA.Vr/d/CWp8be', NOW(), NOW()),
('Charlie Brown', 'charlie@example.com', '$2a$10$oRkxtds7C3lUdM2ZahxVmu5ZDTM.Cqj6e6Tc3gJzA.Vr/d/CWp8be', NOW(), NOW()),
('Diana Prince', 'diana@example.com', '$2a$10$oRkxtds7C3lUdM2ZahxVmu5ZDTM.Cqj6e6Tc3gJzA.Vr/d/CWp8be', NOW(), NOW()),
('Eve Davis', 'eve@example.com', '$2a$10$oRkxtds7C3lUdM2ZahxVmu5ZDTM.Cqj6e6Tc3gJzA.Vr/d/CWp8be', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Insert demo wallets with initial balances
INSERT INTO wallets (user_id, balance, created_at, updated_at)
SELECT id, 
    CASE 
        WHEN email = 'alice@example.com' THEN 1000000.00
        WHEN email = 'bob@example.com' THEN 750000.00
        WHEN email = 'charlie@example.com' THEN 500000.00
        WHEN email = 'diana@example.com' THEN 250000.00
        ELSE 100000.00
    END,
    NOW(),
    NOW()
FROM users
WHERE email IN ('alice@example.com', 'bob@example.com', 'charlie@example.com', 'diana@example.com', 'eve@example.com')
ON CONFLICT (user_id) DO NOTHING;

-- Insert demo top-up transactions
INSERT INTO transactions (sender_id, receiver_id, amount, type, status, created_at, updated_at)
SELECT 
    NULL,
    u.id,
    CASE 
        WHEN u.email = 'alice@example.com' THEN 1000000.00
        WHEN u.email = 'bob@example.com' THEN 750000.00
        WHEN u.email = 'charlie@example.com' THEN 500000.00
        WHEN u.email = 'diana@example.com' THEN 250000.00
        ELSE 100000.00
    END,
    'topup',
    'success',
    NOW() - INTERVAL '7 days',
    NOW() - INTERVAL '7 days'
FROM users u
WHERE u.email IN ('alice@example.com', 'bob@example.com', 'charlie@example.com', 'diana@example.com', 'eve@example.com');

-- Insert demo transfer transactions
-- Alice -> Bob: 50000
INSERT INTO transactions (sender_id, receiver_id, amount, type, status, created_at, updated_at)
SELECT 
    (SELECT id FROM users WHERE email = 'alice@example.com'),
    (SELECT id FROM users WHERE email = 'bob@example.com'),
    50000.00,
    'transfer',
    'success',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '5 days';

-- Bob -> Charlie: 100000
INSERT INTO transactions (sender_id, receiver_id, amount, type, status, created_at, updated_at)
SELECT 
    (SELECT id FROM users WHERE email = 'bob@example.com'),
    (SELECT id FROM users WHERE email = 'charlie@example.com'),
    100000.00,
    'transfer',
    'success',
    NOW() - INTERVAL '4 days',
    NOW() - INTERVAL '4 days';

-- Charlie -> Diana: 75000
INSERT INTO transactions (sender_id, receiver_id, amount, type, status, created_at, updated_at)
SELECT 
    (SELECT id FROM users WHERE email = 'charlie@example.com'),
    (SELECT id FROM users WHERE email = 'diana@example.com'),
    75000.00,
    'transfer',
    'success',
    NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '3 days';

-- Diana -> Eve: 25000
INSERT INTO transactions (sender_id, receiver_id, amount, type, status, created_at, updated_at)
SELECT 
    (SELECT id FROM users WHERE email = 'diana@example.com'),
    (SELECT id FROM users WHERE email = 'eve@example.com'),
    25000.00,
    'transfer',
    'success',
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '2 days';

-- Eve -> Alice: 10000
INSERT INTO transactions (sender_id, receiver_id, amount, type, status, created_at, updated_at)
SELECT 
    (SELECT id FROM users WHERE email = 'eve@example.com'),
    (SELECT id FROM users WHERE email = 'alice@example.com'),
    10000.00,
    'transfer',
    'success',
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '1 day';

-- Update wallet balances to reflect transactions
-- Note: Initial balances already include transaction history
-- This is just for demonstration that final balances are:
-- Alice: 1000000 - 50000 + 10000 = 960000
-- Bob: 750000 + 50000 - 100000 = 700000
-- Charlie: 500000 + 100000 - 75000 = 525000
-- Diana: 250000 + 75000 - 25000 = 300000
-- Eve: 100000 + 25000 - 10000 = 115000

UPDATE wallets SET balance = 960000.00, updated_at = NOW() 
WHERE user_id = (SELECT id FROM users WHERE email = 'alice@example.com');

UPDATE wallets SET balance = 700000.00, updated_at = NOW() 
WHERE user_id = (SELECT id FROM users WHERE email = 'bob@example.com');

UPDATE wallets SET balance = 525000.00, updated_at = NOW() 
WHERE user_id = (SELECT id FROM users WHERE email = 'charlie@example.com');

UPDATE wallets SET balance = 300000.00, updated_at = NOW() 
WHERE user_id = (SELECT id FROM users WHERE email = 'diana@example.com');

UPDATE wallets SET balance = 115000.00, updated_at = NOW() 
WHERE user_id = (SELECT id FROM users WHERE email = 'eve@example.com');

COMMIT;
