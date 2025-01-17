-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT 
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
WITH user_feed_follows AS (
    SELECT * 
    FROM feed_follows
    WHERE feed_follows.user_id = $1
)
SELECT
    users.name AS user_name,
    feeds.name AS feed_name
FROM user_feed_follows
INNER JOIN feeds ON user_feed_follows.feed_id = feeds.id
INNER JOIN users ON user_feed_follows.user_id = users.id;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.feed_id = $1;