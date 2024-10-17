-- name: CreateUser :one 
INSERT INTO users (
    user_id,
    email, 
    password_hash, 
    first_name, 
    last_name,
    phone_number
) VALUES (
   $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;