-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;

-- name: GetFeeds :many
SELECT feeds.created_at, feeds.updated_at, feeds.name, feeds.url, users.name AS user_name FROM feeds
JOIN users
ON users.id = feeds.user_id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) SELECT inserted_feed_follows.*,
    feeds.name as feed_name,
    users.name as user_name
    FROM inserted_feed_follows
JOIN users
ON users.id = inserted_feed_follows.user_id
JOIN feeds
ON feeds.id = inserted_feed_follows.feed_id;

-- name: GetUserFeeds :many
SELECT feed_follows.*,
    users.name AS user_name,
    feeds.name AS feed_name
    FROM feed_follows
JOIN users
ON users.id = feed_follows.user_id
JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND
feed_id = $2;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2,
updated_at = $3
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST,
id ASC
LIMIT 1;
