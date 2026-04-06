-- name: CreateFeed :one
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

-- name: GetFeeds :many
SELECT users.name, feeds.name, feeds.url
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedId :one
SELECT id FROM feeds
WHERE url = $1;


-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
    )
    RETURNING *
) SELECT inserted_feed_follow.*,
      feeds.name as feed_name,
      users.name as user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;