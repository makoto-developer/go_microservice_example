-- Rollback migration: Restore single users table

-- Step 1: Restore users table from backup
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users AS SELECT * FROM users_backup;

-- Restore primary key and constraints
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

-- Restore indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_verification_token ON users(email_verification_token) WHERE email_verification_token IS NOT NULL;
CREATE INDEX idx_users_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;

-- Step 2: Restore refresh_tokens table from backup
DROP TABLE IF EXISTS refresh_tokens CASCADE;
CREATE TABLE refresh_tokens AS SELECT * FROM refresh_tokens_backup;

-- Restore primary key and constraints
ALTER TABLE refresh_tokens ADD PRIMARY KEY (id);
ALTER TABLE refresh_tokens ADD CONSTRAINT refresh_tokens_token_key UNIQUE (token);
ALTER TABLE refresh_tokens ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Restore indexes
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Step 3: Drop new tables
DROP TABLE IF EXISTS customer_refresh_tokens CASCADE;
DROP TABLE IF EXISTS owner_refresh_tokens CASCADE;
DROP TABLE IF EXISTS customer_users CASCADE;
DROP TABLE IF EXISTS owner_users CASCADE;

-- Step 4: Drop backups
DROP TABLE IF EXISTS users_backup;
DROP TABLE IF EXISTS refresh_tokens_backup;

COMMENT ON TABLE users IS 'Restored from backup - single user table with role-based authentication';
