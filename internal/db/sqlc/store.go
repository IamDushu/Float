package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateNurseAccountParams struct {
	Email             string `json:"email"`
	PasswordHash      string `json:"password_hash"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	LicenseNumber     string `json:"license_number"`
	Specialization    string `json:"specialization"`
	YearsOfExperience int32  `json:"years_of_experience"`
	ZipCode           string `json:"zip_code"`
}

type CreateNurseAccountResult struct {
	User  User
	Nurse Nurse
}

func (s *Store) CreateNurseAccountTx(ctx context.Context, arg CreateNurseAccountParams) (CreateNurseAccountResult, error) {
	var result CreateNurseAccountResult

	err := s.execTx(ctx, func(q *Queries) error {
		user, err := q.CreateUser(ctx, CreateUserParams{
			UserID:       uuid.New(),
			Email:        arg.Email,
			PasswordHash: arg.PasswordHash,
			FirstName:    arg.FirstName,
			LastName:     arg.LastName,
			PhoneNumber:  arg.PhoneNumber,
		})
		if err != nil {
			return err
		}
		result.User = user

		nurse, err := q.CreateNurse(ctx, CreateNurseParams{
			NurseID:           uuid.New(),
			UserID:            user.UserID,
			LicenseNumber:     arg.LicenseNumber,
			Specialization:    arg.Specialization,
			YearsOfExperience: arg.YearsOfExperience,
			ZipCode:           arg.ZipCode,
		})
		if err != nil {
			return err
		}
		result.Nurse = nurse

		return nil
	})

	return result, err
}
