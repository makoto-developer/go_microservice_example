-- Auth Service Database Schema
-- 認証サービスのデータベーススキーマ

-- =========================================
-- Enums
-- =========================================

CREATE TYPE user_role AS ENUM (
  'CUSTOMER',
  'SHOP_OWNER',
  'ADMIN'
);

-- =========================================
-- Table: users
-- ユーザー情報
-- =========================================

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  role user_role NOT NULL DEFAULT 'CUSTOMER',
  email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  email_verification_token VARCHAR(255),
  email_verification_expires_at TIMESTAMP WITH TIME ZONE,
  password_reset_token VARCHAR(255),
  password_reset_expires_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  -- 制約
  CONSTRAINT check_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

-- インデックス
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_email_verification_token ON users(email_verification_token) WHERE email_verification_token IS NOT NULL;
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;
CREATE INDEX idx_users_created_at ON users(created_at);

-- =========================================
-- Table: refresh_tokens
-- リフレッシュトークン情報
-- =========================================

CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(512) UNIQUE NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  -- 制約
  CONSTRAINT check_expires_at_future CHECK (expires_at > created_at)
);

-- インデックス
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked) WHERE revoked = FALSE;

-- 有効期限切れのトークンを検索するためのインデックス
CREATE INDEX idx_refresh_tokens_expired ON refresh_tokens(expires_at)
  WHERE revoked = FALSE AND expires_at < NOW();

-- =========================================
-- Functions & Triggers
-- =========================================

-- ユーザー更新時のタイムスタンプ自動更新
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- メール認証完了時にトークンをクリア
CREATE OR REPLACE FUNCTION clear_email_verification_token()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.email_verified = TRUE AND OLD.email_verified = FALSE THEN
    NEW.email_verification_token = NULL;
    NEW.email_verification_expires_at = NULL;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER clear_email_verification_token_on_verify
  BEFORE UPDATE ON users
  FOR EACH ROW
  WHEN (NEW.email_verified = TRUE)
  EXECUTE FUNCTION clear_email_verification_token();

-- パスワードリセット完了時にトークンをクリア
CREATE OR REPLACE FUNCTION clear_password_reset_token()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.password_hash != OLD.password_hash AND OLD.password_reset_token IS NOT NULL THEN
    NEW.password_reset_token = NULL;
    NEW.password_reset_expires_at = NULL;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER clear_password_reset_token_on_reset
  BEFORE UPDATE ON users
  FOR EACH ROW
  WHEN (NEW.password_hash IS DISTINCT FROM OLD.password_hash)
  EXECUTE FUNCTION clear_password_reset_token();

-- 期限切れのリフレッシュトークンを自動削除（定期実行用）
CREATE OR REPLACE FUNCTION cleanup_expired_refresh_tokens()
RETURNS INTEGER AS $$
DECLARE
  deleted_count INTEGER;
BEGIN
  DELETE FROM refresh_tokens
  WHERE expires_at < NOW() - INTERVAL '7 days';
  
  GET DIAGNOSTICS deleted_count = ROW_COUNT;
  RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- Initial Data (Development)
-- =========================================

-- 開発環境用のテストユーザー
-- パスワード: Test1234!
INSERT INTO users (id, email, password_hash, role, email_verified)
VALUES
  ('00000000-0000-0000-0000-000000000001', 'customer@example.com', '$2a$10$rOvHKZqN8K8wZJqmYqYQqOYqN8K8wZJqmYqYQqOYqN8K8wZJqmYqY', 'CUSTOMER', TRUE),
  ('00000000-0000-0000-0000-000000000002', 'shopowner@example.com', '$2a$10$rOvHKZqN8K8wZJqmYqYQqOYqN8K8wZJqmYqYQqOYqN8K8wZJqmYqY', 'SHOP_OWNER', TRUE),
  ('00000000-0000-0000-0000-000000000003', 'admin@example.com', '$2a$10$rOvHKZqN8K8wZJqmYqYQqOYqN8K8wZJqmYqYQqOYqN8K8wZJqmYqY', 'ADMIN', TRUE)
ON CONFLICT (email) DO NOTHING;

-- =========================================
-- Comments
-- =========================================

COMMENT ON TABLE users IS 'ユーザー情報';
COMMENT ON COLUMN users.email IS 'メールアドレス（ログインID）';
COMMENT ON COLUMN users.password_hash IS 'パスワードハッシュ（bcrypt）';
COMMENT ON COLUMN users.role IS 'ユーザーロール';
COMMENT ON COLUMN users.email_verified IS 'メール認証済みフラグ';
COMMENT ON COLUMN users.email_verification_token IS 'メール認証用トークン';
COMMENT ON COLUMN users.email_verification_expires_at IS 'メール認証トークン有効期限';
COMMENT ON COLUMN users.password_reset_token IS 'パスワードリセット用トークン';
COMMENT ON COLUMN users.password_reset_expires_at IS 'パスワードリセットトークン有効期限';

COMMENT ON TABLE refresh_tokens IS 'リフレッシュトークン情報';
COMMENT ON COLUMN refresh_tokens.user_id IS 'ユーザーID';
COMMENT ON COLUMN refresh_tokens.token IS 'リフレッシュトークン（ハッシュ化）';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'トークン有効期限';
COMMENT ON COLUMN refresh_tokens.revoked IS '無効化フラグ';
