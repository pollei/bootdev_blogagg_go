-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetFeedsByName :many
SELECT * FROM feeds
WHERE name = $1 LIMIT 250;

-- name: GetFeedsByUrl :many
SELECT * FROM feeds
WHERE url = $1 LIMIT 25;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1 LIMIT 1;

-- name: GetFeedsByUserId :many
SELECT * FROM feeds
WHERE user_id = $1 LIMIT 250;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsSummary :many
SELECT feeds.name, feeds.url, users.name as user_name
    FROM feeds JOIN users ON feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds SET updated_at = $2 WHERE id = $1 ;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY updated_at ASC NULLS FIRST LIMIT 1;


