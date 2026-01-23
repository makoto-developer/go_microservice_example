-- Create payment_methods table
CREATE TABLE IF NOT EXISTS payment_methods (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    stripe_payment_method_id VARCHAR(255) NOT NULL,
    card_last4 VARCHAR(4) NOT NULL,
    card_brand VARCHAR(50) NOT NULL,
    card_exp_month INTEGER NOT NULL CHECK (card_exp_month >= 1 AND card_exp_month <= 12),
    card_exp_year INTEGER NOT NULL,
    cardholder_name VARCHAR(200) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_methods_customer_id ON payment_methods(customer_id);
CREATE INDEX idx_payment_methods_stripe_id ON payment_methods(stripe_payment_method_id);
CREATE INDEX idx_payment_methods_is_default ON payment_methods(customer_id, is_default);
