-- Seed demo users with fixed UUIDs for consistency
INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES
  ('11111111-1111-1111-1111-111111111111', 'alice', 'alice@example.com', crypt('password123', gen_salt('bf')), NOW(), NOW()),
  ('22222222-2222-2222-2222-222222222222', 'bob', 'bob@example.com', crypt('password123', gen_salt('bf')), NOW(), NOW()),
  ('33333333-3333-3333-3333-333333333333', 'charlie', 'charlie@example.com', crypt('password123', gen_salt('bf')), NOW(), NOW()),
  ('44444444-4444-4444-4444-444444444444', 'diana', 'diana@example.com', crypt('password123', gen_salt('bf')), NOW(), NOW()),
  ('55555555-5555-5555-5555-555555555555', 'eve', 'eve@example.com', crypt('password123', gen_salt('bf')), NOW(), NOW());

-- Seed wallets for the demo users
INSERT INTO wallets (id, user_id, balance, created_at, updated_at)
SELECT gen_random_uuid(), id, 0, NOW(), NOW() FROM users;
