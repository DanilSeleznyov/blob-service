CREATE TABLE IF NOT EXISTS blobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    extension VARCHAR(100) NOT NULL,
    data BYTEA NOT NULL,
    size BIGINT NOT NULL, 
    md5_hash VARCHAR(255) NOT NULL,
    sha256_hash VARCHAR(300),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    category VARCHAR(100),
    tags VARCHAR(50)[],
    metadata JSONB,
    owner_id VARCHAR(255),
    access_level VARCHAR(100)
)