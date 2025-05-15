-- GetAllActiveEvents
-- name: GetAllActiveEventsWithLikes :many
SELECT
    e.*,
    COUNT(l.user_id) AS like_count
FROM pre_go_events e
LEFT JOIN pre_event_like_user l ON e.id = l.event_id
WHERE e.event_active = TRUE
GROUP BY e.id
ORDER BY e.created_at DESC
    LIMIT $1 OFFSET $2;

-- GetEventById
-- name: GetEventById :one
SELECT * FROM pre_go_events WHERE id = $1;

-- AddNewEvent
-- name: AddNewEvent :one
INSERT INTO pre_go_events (
    event_name,
    event_description,
    event_image_url,
    event_active,
    event_start,
    event_end,
    user_id,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
)
RETURNING *;

-- UpdateEvent
-- name: UpdateEvent :exec
UPDATE pre_go_events
SET
    event_name = $2,
    event_description = $3,
    event_image_url = $4,
    event_active = $5,
    event_start = $6,
    event_end = $7,
    user_id = $8,
    updated_at = NOW()
WHERE id = $1;

-- DeleteEvent
-- name: DeleteEvent :exec
DELETE FROM pre_go_events WHERE id = $1;

-- GetEventByName
-- name: GetEventByName :one
SELECT * FROM pre_go_events WHERE event_name = $1;
