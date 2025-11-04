-- name: AddPost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
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

-- name: GetPostsByFeed :many
SELECT posts.*, feeds.name FROM posts
JOIN feeds
ON feeds.id = posts.feed_id
WHERE feed_id = $1
ORDER BY published_at ASC
LIMIT $2;
