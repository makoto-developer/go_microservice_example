CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    logo_url VARCHAR(512),
    owner_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    business_hours VARCHAR(255) NOT NULL,
    return_policy TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_approval',
    published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shops_owner_id ON shops(owner_id);
CREATE INDEX idx_shops_status ON shops(status);
CREATE INDEX idx_shops_published ON shops(published);
