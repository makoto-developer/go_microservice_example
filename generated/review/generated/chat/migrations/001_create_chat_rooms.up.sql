CREATE TABLE IF NOT EXISTS chat_rooms (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    shop_id UUID NOT NULL,
    last_message TEXT,
    last_message_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_customer_shop UNIQUE (customer_id, shop_id)
);

CREATE INDEX idx_chat_rooms_customer_id ON chat_rooms(customer_id);
CREATE INDEX idx_chat_rooms_shop_id ON chat_rooms(shop_id);
