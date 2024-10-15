-- name: CreateUser :one 
INSERT INTO users (
    email, 
    password_hash, 
    first_name, 
    last_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;