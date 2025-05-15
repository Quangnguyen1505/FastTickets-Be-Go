-- Create user like event
-- name: CreateEventLike :one
INSERT INTO pre_event_like_user (event_id, user_id)
VALUES ($1, $2)
    RETURNING event_id, user_id, liked_at;

-- Delete user like event
-- name: DeleteEventLike :exec
DELETE FROM pre_event_like_user
WHERE event_id = $1 AND user_id = $2;

-- Get event user like
-- name: GetEventsUserLike :many
SELECT event_id, user_id
FROM pre_event_like_user
WHERE user_id = $1;
