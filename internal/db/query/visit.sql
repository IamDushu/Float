-- name: CreateVisit :one
INSERT INTO visits (
  visit_id,
  nurse_id,
  patient_id,
  scheduled_at, 
  completed_at, 
  status, 
  notes
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
) RETURNING *;
