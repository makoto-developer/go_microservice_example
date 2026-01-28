-- Admin Service Database Schema
-- Based on: mps-workspace/solutions/admin-service/service.model

-- System Settings Table
CREATE TABLE IF NOT EXISTS system_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    setting_key VARCHAR(255) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    setting_type VARCHAR(50) NOT NULL CHECK (setting_type IN ('STRING', 'NUMBER', 'BOOLEAN', 'JSON')),
    description TEXT,
    updated_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_system_settings_key ON system_settings(setting_key);

-- Categories Table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    parent_id UUID,
    level INTEGER NOT NULL DEFAULT 0,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_categories_level ON categories(level);

-- Forbidden Words Table
CREATE TABLE IF NOT EXISTS forbidden_words (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    word VARCHAR(255) NOT NULL,
    context VARCHAR(50) NOT NULL CHECK (context IN ('REVIEW', 'CHAT', 'ALL')),
    severity VARCHAR(50) NOT NULL CHECK (severity IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forbidden_words_context ON forbidden_words(context);
CREATE INDEX idx_forbidden_words_word ON forbidden_words(word);

-- Audit Logs Table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operation_type VARCHAR(100) NOT NULL,
    operator_id UUID NOT NULL,
    operator_name VARCHAR(255) NOT NULL,
    target_type VARCHAR(100) NOT NULL,
    target_id UUID NOT NULL,
    operation_detail JSONB NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_operator ON audit_logs(operator_id);
CREATE INDEX idx_audit_logs_target ON audit_logs(target_type, target_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_operation_type ON audit_logs(operation_type);

-- Dashboard Metrics Table
CREATE TABLE IF NOT EXISTS dashboard_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metric_type VARCHAR(100) NOT NULL,
    metric_value DECIMAL(20, 2) NOT NULL,
    metric_date DATE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_dashboard_metrics_type_date ON dashboard_metrics(metric_type, metric_date);
CREATE INDEX idx_dashboard_metrics_date ON dashboard_metrics(metric_date DESC);

-- Service Health Check Table
CREATE TABLE IF NOT EXISTS service_health_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('HEALTHY', 'DEGRADED', 'UNHEALTHY')),
    response_time_ms INTEGER NOT NULL,
    error_message TEXT,
    checked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_service_health_service ON service_health_checks(service_name);
CREATE INDEX idx_service_health_checked_at ON service_health_checks(checked_at DESC);

-- Insert default system settings
INSERT INTO system_settings (setting_key, setting_value, setting_type, description, updated_by) VALUES
('maintenance_mode', 'false', 'BOOLEAN', 'システムメンテナンスモード', '00000000-0000-0000-0000-000000000000'),
('max_upload_size_mb', '10', 'NUMBER', '最大アップロードサイズ (MB)', '00000000-0000-0000-0000-000000000000'),
('session_timeout_minutes', '30', 'NUMBER', 'セッションタイムアウト (分)', '00000000-0000-0000-0000-000000000000'),
('shop_approval_required', 'true', 'BOOLEAN', 'ショップ承認必須', '00000000-0000-0000-0000-000000000000')
ON CONFLICT (setting_key) DO NOTHING;
