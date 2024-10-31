package db

import (
	"context"

	"github.com/IamDushu/Float/internal/util"
	"github.com/google/uuid"
)

func (s *Store) ManifestTokenTx(ctx context.Context, verificationRecord EmailVerification) error {
	err := s.execTx(ctx, func(q *Queries) error {
		_, err := q.UpdateVerifyRecordInvalid(ctx, verificationRecord.VerificationID)
		if err != nil {
			return err
		}

		if verificationRecord.Purpose == util.SIGNUP {
			_, err = q.CreateUser(ctx, CreateUserParams{
				UserID: uuid.New(),
				Email:  verificationRecord.Email,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
