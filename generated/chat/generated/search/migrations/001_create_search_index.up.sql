CREATE TYPE index_type AS ENUM ('product', 'shop');

CREATE TABLE IF NOT EXISTS search_index (
    id UUID PRIMARY KEY,
    entity_type index_type NOT NULL,
    entity_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    keywords TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_entity UNIQUE (entity_type, entity_id)
);

CREATE INDEX idx_search_index_entity_type ON search_index(entity_type);
CREATE INDEX idx_search_index_title ON search_index USING gin(to_tsvector('english', title));
CREATE INDEX idx_search_index_description ON search_index USING gin(to_tsvector('english', description));
CREATE INDEX idx_search_index_keywords ON search_index USING gin(to_tsvector('english', keywords));
