-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL, -- Assuming author is a user
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    -- Add a foreign key constraint to the users table
    CONSTRAINT fk_author
        FOREIGN KEY(author_id)
        REFERENCES users(id)
        ON DELETE CASCADE -- If a user is deleted, their articles are also deleted
);

-- Add indexes for faster lookups
CREATE INDEX idx_articles_author_id ON articles (author_id);
CREATE INDEX idx_articles_created_at ON articles (created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd
