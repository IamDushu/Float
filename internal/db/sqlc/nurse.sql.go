// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: nurse.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createNurse = `-- name: CreateNurse :one
INSERT INTO nurses (
  nurse_id,
  user_id,
  license_number,
  specialization,
  years_of_experience,
  zip_code
) VALUES (
   $1, $2, $3, $4, $5, $6
) RETURNING nurse_id, user_id, license_number, specialization, years_of_experience, zip_code, created_at
`

type CreateNurseParams struct {
	NurseID           uuid.UUID `json:"nurse_id"`
	UserID            uuid.UUID `json:"user_id"`
	LicenseNumber     string    `json:"license_number"`
	Specialization    string    `json:"specialization"`
	YearsOfExperience int32     `json:"years_of_experience"`
	ZipCode           string    `json:"zip_code"`
}

func (q *Queries) CreateNurse(ctx context.Context, arg CreateNurseParams) (Nurse, error) {
	row := q.db.QueryRowContext(ctx, createNurse,
		arg.NurseID,
		arg.UserID,
		arg.LicenseNumber,
		arg.Specialization,
		arg.YearsOfExperience,
		arg.ZipCode,
	)
	var i Nurse
	err := row.Scan(
		&i.NurseID,
		&i.UserID,
		&i.LicenseNumber,
		&i.Specialization,
		&i.YearsOfExperience,
		&i.ZipCode,
		&i.CreatedAt,
	)
	return i, err
}

const getNurse = `-- name: GetNurse :one
SELECT nurse_id, user_id, license_number, specialization, years_of_experience, zip_code, created_at FROM nurses
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetNurse(ctx context.Context, userID uuid.UUID) (Nurse, error) {
	row := q.db.QueryRowContext(ctx, getNurse, userID)
	var i Nurse
	err := row.Scan(
		&i.NurseID,
		&i.UserID,
		&i.LicenseNumber,
		&i.Specialization,
		&i.YearsOfExperience,
		&i.ZipCode,
		&i.CreatedAt,
	)
	return i, err
}