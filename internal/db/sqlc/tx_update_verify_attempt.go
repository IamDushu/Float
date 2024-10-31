package db

import (
	"context"

	"github.com/google/uuid"
)

func (s *Store) UpdateVerifyAttemptTx(ctx context.Context, verificationID uuid.UUID) (EmailVerification, error) {
	var verifyRecord EmailVerification

	err := s.execTx(ctx, func(q *Queries) error {
		record, err := q.UpdateVerifyRecordAttempt(ctx, verificationID)
		if err != nil {
			return err
		}

		if record.Attempts == 5 {
			record, err = q.UpdateVerifyRecordInvalid(ctx, verificationID)
			if err != nil {
				return err
			}
		}

		verifyRecord = record
		return nil
	})

	return verifyRecord, err
}
