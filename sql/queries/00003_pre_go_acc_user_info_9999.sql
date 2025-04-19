-- name: GetUser :one
SELECT 
    user_id,
    user_account,
    user_nickname, 
    user_avatar, 
    user_state, 
    user_mobile, 
    user_gender, 
    user_birthday, 
    user_email, 
    user_is_authentication,
    created_at,
    updated_at
FROM pre_go_acc_user_info_9999
WHERE user_id = $1;

-- name: GetUsers :many
SELECT 
    user_id,
    user_account,
    user_nickname, 
    user_avatar, 
    user_state, 
    user_mobile, 
    user_gender, 
    user_birthday, 
    user_email, 
    user_is_authentication,
    created_at,
    updated_at
FROM pre_go_acc_user_info_9999
WHERE user_id = ANY($1::bigint[]);

-- name: findUsers :many
SELECT *
FROM pre_go_acc_user_info_9999
WHERE user_account LIKE '%' || $1 || '%' OR user_email LIKE '%' || $2 || '%';

-- name: listUsers :many
SELECT *
FROM pre_go_acc_user_info_9999
ORDER BY created_at DESC;

-- name: RemoveUser :exec
DELETE FROM pre_go_acc_user_info_9999 WHERE user_id = $1;

-- name: UpdatePassword :exec


-- name: AddUserAutoUserId :execresult
INSERT INTO pre_go_acc_user_info_9999 (
    user_account, user_nickname, user_avatar, user_state, user_mobile,
    user_gender, user_birthday, user_email, user_is_authentication
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: AddUserHaveUserId :one
INSERT INTO pre_go_acc_user_info_9999 (
    user_id, user_account, user_nickname, user_avatar, user_state, user_mobile,
    user_gender, user_birthday, user_email, user_is_authentication
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING user_id;
-- name: EditUserByUserId :execresult
UPDATE pre_go_acc_user_info_9999
SET user_nickname = $1, user_avatar = $2, user_mobile = $3,
    user_gender = $4, user_birthday = $5, user_email = $6, updated_at = NOW()
WHERE user_id = $7 AND user_is_authentication = 1;



