-- name: CreateNurse :one
INSERT INTO nurses (
  nurse_id,
  user_id,
  license_number,
  specialization,
  years_of_experience,
  zip_code
) VALUES (
   $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetNurse :one
SELECT * FROM nurses
WHERE user_id = $1 LIMIT 1;