-- name: CreateContactMessage :one
INSERT INTO pre_go_contact_messages (
    name, email, phone, message, status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetContactMessageByID :one
SELECT * FROM pre_go_contact_messages
WHERE id = $1;

-- name: GetAllContactMessages :many
SELECT id, name, email, phone, message, status, created_at
FROM pre_go_contact_messages
WHERE ($1 = -1 OR status = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: UpdateContactMessageStatus :exec
UPDATE pre_go_contact_messages
SET status = $2
WHERE id = $1;

-- name: DeleteContactMessage :exec
DELETE FROM pre_go_contact_messages
WHERE id = $1;
