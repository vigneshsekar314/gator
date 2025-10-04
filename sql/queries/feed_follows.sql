
-- name: CreateFeedFollow :one
WITH A AS (INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ( $1, $2, $3, $4, $5) RETURNING *) 
 SELECT A.id AS feed_follows_id, A.created_at, A.updated_at, A.user_id, A.feed_id, users.name AS username, feeds.name AS feed_name FROM A JOIN users ON A.user_id = users.id JOIN feeds ON A.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT users.name AS user_name, feeds.name AS feed_name  FROM users JOIN feed_follows ON users.id = feed_follows.user_id
JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE users.id = $1;
