-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY, -- https://www.postgresql.org/docs/16/datatype-uuid.html
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT, -- the title of the post
    url TEXT NOT NULL, -- the URL of the post (this should be unique)
    description TEXT, -- the description of the post
    published_at TIMESTAMP NOT NULL, -- the time the post was published
    feed_id UUID NOT NULL
        REFERENCES feeds (id) ON DELETE CASCADE,
    UNIQUE(  url )
);  -- https://www.postgresql.org/docs/16/sql-createtable.html

-- +goose Down
DROP TABLE posts;