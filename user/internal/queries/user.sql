-- name: GetUser :one
SELECT 
    users.uid, users.email, users.username, users.secret_key,
    users.first_name, users.last_name, 
    users.phone_number, users.is_verified, users.last_login,
    user_otps.otp_enabled, user_otps.otp_verified,
    user_otps.otp_url, user_otps.otp_secret
FROM users
LEFT JOIN user_otps
ON users.uid = user_otps.uid
WHERE (users.uid = $1 OR users.phone_number = $1 ) AND users.is_active = TRUE LIMIT 1;

-- name: UpdateUser :one
UPDATE users SET
    uid = $1,
    phone_number = $2,
    email = $3,
    username = $4,
    first_name = $5,
    last_name = $6
WHERE "uid" = $1
RETURNING *;

-- name: UpdateUserVerified :exec
UPDATE users SET 
    is_verified = $1 
WHERE "uid" = $2;

-- name: DeleteUser :exec
UPDATE users SET
    is_active = FALSE
WHERE "uid" = $1;