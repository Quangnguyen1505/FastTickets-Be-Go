-- file: pre_go_acc_user_two_factor_9999.sql

-- EnableTwoFactor
-- name: EnableTwoFactorTypeEmail :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (
    user_id,
    two_factor_auth_type,
    two_factor_auth_secret,
    two_factor_email,
    two_factor_is_active, 
    two_factor_created_at, 
    two_factor_updated_at
) VALUES (
    $1, $2, 'OTP', $3, FALSE, NOW(), NOW()
)
RETURNING *;

-- DisableTwoFactor
-- name: DisableTwoFactor :exec
UPDATE pre_go_acc_user_two_factor_9999
  SET two_factor_is_active = FALSE,
  two_factor_updated_at = NOW()
WHERE user_id = $1 AND two_factor_auth_type = $2;

-- UpdateTwoFactorStatusVerification
-- name: UpdateTwoFactorStatusVerification :exec
UPDATE pre_go_acc_user_two_factor_9999
  SET two_factor_is_active = TRUE,
  two_factor_updated_at = NOW()
WHERE user_id = $1 AND two_factor_auth_type = $2 AND two_factor_is_active = FALSE;

-- VerifyTwoFactor
-- name: VerifyTwoFactor :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = $1 AND two_factor_auth_type = $2 AND two_factor_is_active = TRUE;

-- GetTwoFactorStatus
-- name: GetTwoFactorStatus :one
SELECT two_factor_is_active
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = $1 AND two_factor_auth_type = $2;

-- IsTwoFactorEnabled
-- name: IsTwoFactorEnabled :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = $1 AND two_factor_is_active = TRUE;

-- AddOrUpdatePhoneNumber
-- name: AddOrUpdatePhoneNumber :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (
    user_id,
    two_factor_phone,
    two_factor_is_active
) VALUES (
    $1, $2, TRUE
) ON CONFLICT (user_id) DO UPDATE SET
    two_factor_phone = EXCLUDED.two_factor_phone,
    two_factor_updated_at = NOW();

-- AddOrUpdateEmail
-- name: AddOrUpdateEmail :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (
    user_id,
    two_factor_email,
    two_factor_is_active
) VALUES (
    $1, $2, TRUE
) ON CONFLICT (user_id) DO UPDATE SET
    two_factor_email = EXCLUDED.two_factor_email,
    two_factor_updated_at = NOW();

-- GetUserTwoFactorMethods
-- name: GetUserTwoFactorMethods :many
SELECT *
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = $1;

-- ReactivateTwoFactor
-- name: ReactivateTwoFactor :exec
UPDATE pre_go_acc_user_two_factor_9999
  SET two_factor_is_active = TRUE,
  two_factor_updated_at = NOW()
WHERE user_id = $1 AND two_factor_auth_type = $2;

-- RemoveTwoFactor
-- name: RemoveTwoFactor :exec
DELETE FROM pre_go_acc_user_two_factor_9999
WHERE user_id = $1 AND two_factor_auth_type = $2;

-- CountActiveTwoFactorMethods
-- name: CountActiveTwoFactorMethods :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = $1 AND two_factor_is_active = TRUE;

-- GetTwoFactorMethodByID
-- name: GetTwoFactorMethodByID :one
SELECT *
FROM pre_go_acc_user_two_factor_9999
WHERE two_factor_id = $1;

-- GetTwoFactorMethodByIDandType
-- name: GetTwoFactorMethodByIDandType :one
SELECT *
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = $1 AND two_factor_auth_type = $2;