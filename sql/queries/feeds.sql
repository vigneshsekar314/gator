-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
  ) RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsWithUsername :many
 SELECT feeds.*, users.name AS username FROM feeds JOIN feed_follows ON feeds.id = feed_follows.feed_id JOIN users ON feed_follows.user_id = users.id;

-- name: GetFeedsByUrl :one
SELECT feeds.id, feeds.name FROM feeds WHERE feeds.url = $1 LIMIT 1;

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT feeds.id, feeds.name, feeds.url FROM feeds ORDER BY last_fetched_at DESC NULLS FIRST LIMIT 1;
