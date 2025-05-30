-- First, create ENUM types for status and category
CREATE TYPE campaign_status AS ENUM (
    'active',
    'paused',
    'completed',
    'cancelled'
);

CREATE TYPE campaign_category AS ENUM (
    'education',
    'healthcare',
    'environment',
    'animals',
    'emergency',
    'community',
    'technology',
    'arts',
    'sports'
);

ALTER TYPE campaign_category ADD VALUE 'unspecified' BEFORE 'education';

-- Now, create the campaigns table
CREATE TABLE campaigns (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(255),
    target_amount INTEGER NOT NULL,
    collected_amount INTEGER DEFAULT 0,
    deadline DATE NOT NULL,
    status campaign_status NOT NULL DEFAULT 'active',
    category campaign_category NOT NULL,
    min_donation INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
