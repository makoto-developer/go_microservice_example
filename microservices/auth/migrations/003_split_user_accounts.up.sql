-- Migration: Split user accounts into customer_users and owner_users
-- Description: Separate customer and owner authentication (Amazon-style)

-- Step 1: Create customer_users table
CREATE TABLE IF NOT EXISTS customer_users (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                           VARCHAR(255) NOT NULL UNIQUE,
    password_hash                   VARCHAR(255) NOT NULL,
    email_verified                  BOOLEAN NOT NULL DEFAULT FALSE,
    email_verification_token        VARCHAR(255),
    email_verification_expires_at   TIMESTAMP,
    password_reset_token            VARCHAR(255),
    password_reset_expires_at       TIMESTAMP,
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customer_users_email ON customer_users(email);
CREATE INDEX idx_customer_users_verification_token ON customer_users(email_verification_token) WHERE email_verification_token IS NOT NULL;
CREATE INDEX idx_customer_users_reset_token ON customer_users(password_reset_token) WHERE password_reset_token IS NOT NULL;

COMMENT ON TABLE customer_users IS 'Customer authentication accounts - separate from owner accounts';
COMMENT ON COLUMN customer_users.email IS 'Can overlap with owner_users.email (same person, different account)';

-- Step 2: Create owner_users table
CREATE TABLE IF NOT EXISTS owner_users (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                           VARCHAR(255) NOT NULL UNIQUE,
    password_hash                   VARCHAR(255) NOT NULL,
    email_verified                  BOOLEAN NOT NULL DEFAULT FALSE,
    email_verification_token        VARCHAR(255),
    email_verification_expires_at   TIMESTAMP,
    password_reset_token            VARCHAR(255),
    password_reset_expires_at       TIMESTAMP,
    business_verified               BOOLEAN NOT NULL DEFAULT FALSE,
    business_verification_status    VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_owner_users_email ON owner_users(email);
CREATE INDEX idx_owner_users_verification_token ON owner_users(email_verification_token) WHERE email_verification_token IS NOT NULL;
CREATE INDEX idx_owner_users_reset_token ON owner_users(password_reset_token) WHERE password_reset_token IS NOT NULL;
CREATE INDEX idx_owner_users_business_status ON owner_users(business_verification_status);

COMMENT ON TABLE owner_users IS 'Shop owner authentication accounts - separate from customer accounts';
COMMENT ON COLUMN owner_users.business_verified IS 'Whether the business has been verified by admin';
COMMENT ON COLUMN owner_users.business_verification_status IS 'pending, approved, rejected';

-- Step 3: Migrate existing data from users table
-- Customer accounts (role = 'CUSTOMER')
INSERT INTO customer_users (
    id,
    email,
    password_hash,
    email_verified,
    email_verification_token,
    email_verification_expires_at,
    password_reset_token,
    password_reset_expires_at,
    created_at,
    updated_at
)
SELECT
    id,
    email,
    password_hash,
    email_verified,
    email_verification_token,
    email_verification_expires_at,
    password_reset_token,
    password_reset_expires_at,
    created_at,
    updated_at
FROM users
WHERE role = 'CUSTOMER'
ON CONFLICT (email) DO NOTHING;

-- Owner accounts (role = 'SHOP_OWNER')
INSERT INTO owner_users (
    id,
    email,
    password_hash,
    email_verified,
    email_verification_token,
    email_verification_expires_at,
    password_reset_token,
    password_reset_expires_at,
    business_verified,
    business_verification_status,
    created_at,
    updated_at
)
SELECT
    id,
    email,
    password_hash,
    email_verified,
    email_verification_token,
    email_verification_expires_at,
    password_reset_token,
    password_reset_expires_at,
    FALSE,  -- business_verified default
    'pending',  -- business_verification_status default
    created_at,
    updated_at
FROM users
WHERE role = 'SHOP_OWNER'
ON CONFLICT (email) DO NOTHING;

-- Step 4: Create refresh_tokens tables for each user type
CREATE TABLE IF NOT EXISTS customer_refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES customer_users(id) ON DELETE CASCADE,
    token       VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customer_refresh_tokens_user_id ON customer_refresh_tokens(user_id);
CREATE INDEX idx_customer_refresh_tokens_token ON customer_refresh_tokens(token);
CREATE INDEX idx_customer_refresh_tokens_expires_at ON customer_refresh_tokens(expires_at);

CREATE TABLE IF NOT EXISTS owner_refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES owner_users(id) ON DELETE CASCADE,
    token       VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_owner_refresh_tokens_user_id ON owner_refresh_tokens(user_id);
CREATE INDEX idx_owner_refresh_tokens_token ON owner_refresh_tokens(token);
CREATE INDEX idx_owner_refresh_tokens_expires_at ON owner_refresh_tokens(expires_at);

-- Step 5: Migrate existing refresh_tokens
-- Customer refresh tokens
INSERT INTO customer_refresh_tokens (id, user_id, token, expires_at, created_at)
SELECT rt.id, rt.user_id, rt.token, rt.expires_at, rt.created_at
FROM refresh_tokens rt
INNER JOIN users u ON rt.user_id = u.id
WHERE u.role = 'CUSTOMER'
ON CONFLICT (token) DO NOTHING;

-- Owner refresh tokens
INSERT INTO owner_refresh_tokens (id, user_id, token, expires_at, created_at)
SELECT rt.id, rt.user_id, rt.token, rt.expires_at, rt.created_at
FROM refresh_tokens rt
INNER JOIN users u ON rt.user_id = u.id
WHERE u.role = 'SHOP_OWNER'
ON CONFLICT (token) DO NOTHING;

-- Step 6: Create backup of old tables (optional, can be dropped later)
CREATE TABLE IF NOT EXISTS users_backup AS SELECT * FROM users;
CREATE TABLE IF NOT EXISTS refresh_tokens_backup AS SELECT * FROM refresh_tokens;

COMMENT ON TABLE users_backup IS 'Backup of users table before account split - can be dropped after verification';
COMMENT ON TABLE refresh_tokens_backup IS 'Backup of refresh_tokens table before account split - can be dropped after verification';
