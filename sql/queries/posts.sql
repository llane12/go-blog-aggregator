-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name FROM posts
JOIN feeds ON posts.feed_id = feeds.id
WHERE feed_id IN
(
    SELECT feed_id FROM feed_follows WHERE feed_follows.user_id = $1
)
ORDER BY published_at DESC
LIMIT $2;
