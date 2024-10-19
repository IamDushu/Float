-- name: CreatePatient :one
INSERT INTO patients (
    patient_id, 
    user_id,
    date_of_birth, 
    emergency_contact_name,  
    emergency_contact_phone,  
    medical_history,  
    allergies  
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetPatient :one
SELECT * FROM patients
WHERE user_id = $1 LIMIT 1;