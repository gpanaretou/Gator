-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, url, description, published_at, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT 
    posts.id,
    posts.title,
    posts.url,
    posts.description,
    posts.published_at,
    posts.created_at,
    posts.updated_at,
    feeds.name AS feed_name
FROM posts
INNER JOIN feeds ON posts.feed_id = feeds.id
INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;