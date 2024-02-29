-- name: CreateUser :one
INSERT INTO users (
    uid,
    phone_number,
    email,
    username,
    first_name,
    last_name,
    password,
    is_active,
    is_verified,
    secret_key,
    last_login
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: CreateUserOTP :exec
INSERT INTO user_otps (
    uid,
    otp_enabled,
    otp_verified,
    otp_secret,
    otp_url
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetUserByEmailOrUsername :one
SELECT * FROM users WHERE username = $1 OR email = $1 AND is_active = TRUE;

-- name: GetUserByEmailAndUsername :one
SELECT * FROM users WHERE username = $1 AND email = $2;

-- name: GetUserByUid :one
SELECT * FROM users WHERE uid = $1 LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO "sessions" (
    id,
    uid,
    refresh_token,
    user_agent,
    client_ip,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM "sessions" WHERE id = $1;

-- name: CreateRecoveryAccount :one
INSERT INTO recover_accounts (
    id,
    uid,
    email,
    secret_key,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetRecoverAccount :one
SELECT * FROM recover_accounts WHERE uid = $1 AND secret_key = $2;

-- name: UpdateRecoverAccount :exec
UPDATE recover_accounts
SET 
    is_used = TRUE
WHERE 
    uid = @uid 
    AND secret_key = @secret_key
    AND is_used = FALSE;

-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    id,
    uid,
    email,
    secret_key,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM verify_emails WHERE uid = $1 AND secret_key = $2;

-- name: UpdateVerifyEmail :exec
UPDATE verify_emails
SET 
    is_used = TRUE
WHERE 
    uid = @uid 
    AND secret_key = @secret_key
    AND is_used = FALSE;

-- name: DeleteSessionUser :exec
DELETE FROM "sessions" WHERE uid = $1;

-- name: UpdateUserVerified :exec
UPDATE users SET 
    is_verified = $1 
WHERE "uid" = $2;

-- name: UpdateUserPassword :exec
UPDATE users SET 
    password = $1
WHERE "uid" = $2;

-- name: UpdateUserLastLogin :exec
UPDATE users SET
    last_login = 'now()'
WHERE "uid" = $1;

-- name: UpdateUserActive :exec
UPDATE users SET
    is_active = $1
WHERE "uid" = $2;