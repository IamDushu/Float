-- name: CreateVerifyRecord :one
INSERT INTO email_verification (
    verification_id,
    email,
    token,
    hashed_otp,
    purpose,
    attempts,
    expires_at,
    valid
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetVerifyRecord :one
SELECT * FROM email_verification
WHERE email = $1
AND valid = TRUE
AND purpose = $2;

-- name: GetVerifyRecordOnToken :one
SELECT * FROM email_verification
WHERE token = $1;

-- name: UpdateVerifyRecordInvalid :one
UPDATE email_verification
SET valid = FALSE
WHERE verification_id = $1
RETURNING *;

-- name: UpdateVerifyRecordAttempt :one
UPDATE email_verification
SET attempts = attempts + 1
WHERE verification_id = $1
RETURNING *;