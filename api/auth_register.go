package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/IamDushu/Float/internal/token"
	"github.com/IamDushu/Float/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Mode  string `json:"mode" binding:"required,oneof=signup login"`
}

type registerUserResponse struct {
	Token string `json:"token"`
}

// registerUser handles signup or login requests
func (s *Server) registerUser(ctx *gin.Context) {
	var request registerUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	request.Email = util.NormalizeEmail(request.Email)

	switch request.Mode {
	case util.SIGNUP:
		s.handleSignupUser(ctx, request)
	case util.LOGIN:
		s.handleLoginUser(ctx, request)
	}
}

// handleSignupUser processes signup requests
func (s *Server) handleSignupUser(ctx *gin.Context, req registerUserRequest) {

	var response registerUserResponse

	// Check if user already exists
	if _, err := s.store.GetUser(ctx, req.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User not found, process for new signup
			if err := s.processExistingVerifyRecord(ctx, req.Email, req.Mode); err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			token, otp, err := s.createNewVerifyRecord(ctx, req.Email, req.Mode, true)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			//TODO
			fmt.Println("OTP Sent via email" + strconv.Itoa(otp))

			response.Token = token
			ctx.JSON(http.StatusOK, response)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// User already exists, create non-valid verification record
	token, _, err := s.createNewVerifyRecord(ctx, req.Email, req.Mode, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response.Token = token
	ctx.JSON(http.StatusOK, response)
}

// handleLoginUser processes login requests
func (s *Server) handleLoginUser(ctx *gin.Context, req registerUserRequest) {

	var response registerUserResponse

	// Check if user already exists
	if _, err := s.store.GetUser(ctx, req.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User not found, create non-valid verification record
			token, _, err := s.createNewVerifyRecord(ctx, req.Email, req.Mode, false)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			response.Token = token
			ctx.JSON(http.StatusOK, response)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	// User already exists, process login
	if err := s.processExistingVerifyRecord(ctx, req.Email, req.Mode); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, otp, err := s.createNewVerifyRecord(ctx, req.Email, req.Mode, true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//TODO
	fmt.Println("OTP Sent via email" + strconv.Itoa(otp))

	response.Token = token
	ctx.JSON(http.StatusOK, response)

}

// processExistingVerifyRecord checks and invalidates any existing verification records for a user
func (s *Server) processExistingVerifyRecord(ctx *gin.Context, email string, mode string) error {
	pastRecordArgs := db.GetVerifyRecordParams{
		Email:   email,
		Purpose: mode,
	}

	// Check if a verification record already exists
	pastRecord, err := s.store.GetVerifyRecord(ctx, pastRecordArgs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil // No existing record, nothing to invalidate
		}
		return fmt.Errorf("error retrieving verification record: %w", err)
	}

	// Invalidate existing verification record
	if _, err := s.store.UpdateVerifyRecordInvalid(ctx, pastRecord.VerificationID); err != nil {
		return fmt.Errorf("error invalidating existing verification record: %w", err)
	}
	return nil
}

// createNewVerifyRecord generates a new verification record for a user
func (s *Server) createNewVerifyRecord(ctx *gin.Context, email string, mode string, validity bool) (string, int, error) {
	recordArgs, otp, err := createVerifyRecordParams(email, mode, validity)
	if err != nil {
		return "", 0, fmt.Errorf("error creating verification record params: %w", err)
	}

	record, err := s.store.CreateVerifyRecord(ctx, *recordArgs)
	if err != nil {
		return "", 0, fmt.Errorf("error saving new verification record: %w", err)
	}
	return record.Token, otp, nil
}

func createVerifyRecordParams(email string, purpose string, validity bool) (*db.CreateVerifyRecordParams, int, error) {
	claims := token.Claims{
		Sub: email,
		Iat: time.Now().Unix(),
		Nbf: time.Now().Unix(),
		Exp: time.Now().Add(30 * time.Minute).Unix(),
	}

	tkn, err := token.CreateUnsignedJWT(claims)
	if err != nil {
		return &db.CreateVerifyRecordParams{}, 0, err
	}

	otp, err := util.GenerateOTP()
	if err != nil {
		return &db.CreateVerifyRecordParams{}, 0, err
	}

	hashedOtp, err := util.HashThis(otp)
	if err != nil {
		return &db.CreateVerifyRecordParams{}, 0, err
	}

	verifyRecord := db.CreateVerifyRecordParams{
		VerificationID: uuid.New(),
		Email:          email,
		Token:          tkn,
		HashedOtp:      hashedOtp,
		Purpose:        purpose,
		Attempts:       0,
		ExpiresAt:      time.Unix(claims.Exp, 0),
		Valid:          validity,
	}

	return &verifyRecord, otp, nil
}
