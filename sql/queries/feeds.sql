-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds
ORDER BY user_id;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1 LIMIT 1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET 
    last_fetched_at = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST;

-- name: GetTotalNumberOfFeeds :one
SELECT COUNT(*) FROM feeds;