-- GetAllMenuFunctionActive
-- name: GetAllMenuFunctionActive :many
SELECT * FROM pre_go_menu_function WHERE active = TRUE;

-- GetMenuFunctionById
-- name: GetMenuFunctionById :one
SELECT * FROM pre_go_menu_function WHERE id = $1;

-- AddNewMenuFunction
-- name: AddNewMenuFunction :one
INSERT INTO pre_go_menu_function (
    id, name, description, url, active, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW()
)
RETURNING *;

-- EditMenuFunction
-- name: EditMenuFunction :one
UPDATE pre_go_menu_function
SET name = $2, description = $3, url = $4, active = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- RemoveMenuFunction
-- name: RemoveMenuFunction :one
DELETE FROM pre_go_menu_function WHERE id = $1
RETURNING *;

-- GetMenuFunctionByName
-- name: GetMenuFunctionByName :one
SELECT * FROM pre_go_menu_function WHERE name = $1;

-- GetAllMenuFunctions
-- name: GetAllMenuFunctions :many
SELECT * FROM pre_go_menu_function ORDER BY created_at DESC;