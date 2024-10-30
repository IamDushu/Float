package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/IamDushu/Float/internal/token"
	"github.com/IamDushu/Float/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	SIGNUP = "signup"
	LOGIN  = "login"
)

type registerUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Mode  string `json:"mode" binding:"required,oneof=signup login"`
}

type registerUserResponse struct {
	Token string `json:"token"`
}

func (s *Server) registerUser(ctx *gin.Context) {
	var request registerUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	request.Email = util.NormalizeEmail(request.Email)

	switch request.Mode {
	case SIGNUP:
		s.handleSignupUser(ctx, request)
	case LOGIN:
		ctx.JSON(http.StatusOK, registerUserResponse{Token: "req for login"})
	}
}

func (s *Server) handleSignupUser(ctx *gin.Context, req registerUserRequest) {

	var response registerUserResponse

	_, err := s.store.GetUser(ctx, req.Email)

	if errors.Is(err, sql.ErrNoRows) {
		//User doesn't exist in db
		recordArgs, err := createVerifyRecordParams(req.Email, SIGNUP, true)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("something went wrong while signing up")))
			return
		}
		record, err := s.store.CreateVerifyRecord(ctx, *recordArgs)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		response.Token = record.Token
		ctx.JSON(http.StatusOK, response)
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	recordArgs, err := createVerifyRecordParams(req.Email, SIGNUP, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("something went wrong while signing up")))
		return
	}
	record, err := s.store.CreateVerifyRecord(ctx, *recordArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response.Token = record.Token
	ctx.JSON(http.StatusOK, response)
}

// func (s *Server) handleLoginUser(req registerUserRequest) {

// }

func createVerifyRecordParams(email string, purpose string, validity bool) (*db.CreateVerifyRecordParams, error) {
	claims := token.Claims{
		Sub: email,
		Iat: time.Now(),
		Nbf: time.Now(),
		Exp: time.Now().Add(30 * time.Minute),
	}

	tkn, err := token.CreateUnsignedJWT(claims)
	if err != nil {
		return &db.CreateVerifyRecordParams{}, err
	}

	otp, err := util.GenerateOTP()
	if err != nil {
		return &db.CreateVerifyRecordParams{}, err
	}

	hashedOtp, err := util.HashThis(otp)
	if err != nil {
		return &db.CreateVerifyRecordParams{}, err
	}

	verifyRecord := db.CreateVerifyRecordParams{
		VerificationID: uuid.New(),
		Email:          email,
		Token:          tkn,
		HashedOtp:      hashedOtp,
		Purpose:        purpose,
		Attempts:       0,
		ExpiresAt:      claims.Exp,
		Valid:          validity,
	}

	return &verifyRecord, nil
}
