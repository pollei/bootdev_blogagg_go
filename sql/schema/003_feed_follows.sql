-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY, -- https://www.postgresql.org/docs/16/datatype-uuid.html
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL
        REFERENCES users (id) ON DELETE CASCADE,
    feed_id UUID NOT NULL
        REFERENCES feeds (id) ON DELETE CASCADE,
    UNIQUE(  user_id, feed_id )
);  -- https://www.postgresql.org/docs/16/sql-createtable.html

-- +goose Down
DROP TABLE feed_follows;