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
 SELECT feeds.*, users.name AS username FROM feeds JOIN feed_follows ON feeds.feed_id = feed_follows.id JOIN users ON feed_follows.user_id = users.id;

-- name: GetFeedsByUrl :one
SELECT feeds.id, feeds.name FROM feeds WHERE feeds.url = $1 LIMIT 1;
