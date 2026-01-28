-- Notification Service Database Schema
-- This schema supports notification template management and notification sending

CREATE TABLE IF NOT EXISTS notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    channel VARCHAR(50) NOT NULL,
    subject VARCHAR(255),
    body_template TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE notification_templates IS 'Template definitions for various notification types';
COMMENT ON COLUMN notification_templates.name IS 'Unique template name (e.g., order_confirmation, payment_success)';
COMMENT ON COLUMN notification_templates.channel IS 'Notification channel: email, sms, push';
COMMENT ON COLUMN notification_templates.subject IS 'Email subject (null for SMS/push)';
COMMENT ON COLUMN notification_templates.body_template IS 'Template body with variable placeholders';

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID REFERENCES notification_templates(id),
    recipient VARCHAR(255) NOT NULL,
    channel VARCHAR(50) NOT NULL,
    subject VARCHAR(255),
    body TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE notifications IS 'Notification records for tracking sent and pending notifications';
COMMENT ON COLUMN notifications.status IS 'Status: pending, sent, failed, retrying';
COMMENT ON COLUMN notifications.sent_at IS 'Timestamp when notification was successfully sent';
COMMENT ON COLUMN notifications.error_message IS 'Error message if sending failed';

-- Indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient ON notifications(recipient);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_template_id ON notifications(template_id);

-- Insert sample notification templates
INSERT INTO notification_templates (name, channel, subject, body_template) VALUES
    ('order_confirmation', 'email', 'Order Confirmation - Order #{order_number}',
     'Dear {customer_name},\n\nYour order #{order_number} has been confirmed.\n\nOrder Details:\n{order_details}\n\nTotal Amount: {total_amount}\n\nThank you for your purchase!'),
    ('payment_success', 'email', 'Payment Successful - Order #{order_number}',
     'Dear {customer_name},\n\nYour payment for order #{order_number} has been successfully processed.\n\nAmount Paid: {amount}\nPayment Method: {payment_method}\n\nThank you!'),
    ('shipping_update', 'email', 'Shipping Update - Order #{order_number}',
     'Dear {customer_name},\n\nYour order #{order_number} has been shipped.\n\nTracking Number: {tracking_number}\nEstimated Delivery: {estimated_delivery}\n\nTrack your package here: {tracking_url}'),
    ('order_cancelled', 'email', 'Order Cancelled - Order #{order_number}',
     'Dear {customer_name},\n\nYour order #{order_number} has been cancelled.\n\nReason: {cancellation_reason}\n\nIf you have any questions, please contact our support team.')
ON CONFLICT (name) DO NOTHING;
