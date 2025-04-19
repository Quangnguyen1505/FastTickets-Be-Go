-- name: GetValidOtp :one
SELECT verify_otp, verify_key_hash, verify_key, verify_id 
FROM pre_go_acc_user_verify_9999
WHERE verify_key_hash = $1 AND is_verified = 0;

-- update lai
-- name: UpdateUserValificationStatus :exec
UPDATE pre_go_acc_user_verify_9999
  set is_verified = 1,
  verify_updated_at = now()
WHERE verify_key_hash = $1;

-- name: InsertOtpVerify :one
INSERT INTO pre_go_acc_user_verify_9999 (
  verify_otp, 
  verify_key, 
  verify_key_hash, 
  verify_type,
  is_verified,
  is_deleted,
  verify_created_at,
  verify_updated_at
) VALUES (
  $1, $2, $3, $4, 0, 0, NOW(), NOW()
)
RETURNING verify_id;

-- name: GetInfoOtp :one
SELECT verify_id, verify_otp, verify_key, verify_key_hash, verify_type, is_verified, is_deleted, verify_created_at, verify_updated_at
FROM pre_go_acc_user_verify_9999
WHERE verify_key_hash = $1;
