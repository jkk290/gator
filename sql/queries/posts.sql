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

-- name: GetPostsByFeeds :many
SELECT * FROM posts
JOIN feeds
ON feeds.id = posts.feed_id
WHERE feed_id = ANY($1::uuid[])
ORDER BY published_at DESC
LIMIT $2;
