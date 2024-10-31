package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/IamDushu/Float/internal/util"
	"github.com/gin-gonic/gin"
)

type verifyUserRequest struct {
	Token  string `json:"token" binding:"required,jwt"`
	Digits string `json:"digits" binding:"required,len=5"`
}

type verifyUserResponse struct {
	AccessToken string `json:"access_token"`
	Mode        string `json:"mode"`
	Email       string `json:"email"`
}

func (s *Server) verifyUser(ctx *gin.Context) {
	var request verifyUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	verifyRecord, err := s.store.GetVerifyRecordOnToken(ctx, request.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("something went wrong with the link you used. please go back and try again. [invalid_jwt]")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !verifyRecord.Valid || time.Now().After(verifyRecord.ExpiresAt) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("this code has been used or expired. please go back to get a new code. [used_or_expired]")))
		return
	}

	if err := util.HashVerify(request.Digits, verifyRecord.HashedOtp); err != nil {
		// Updates attempt +1 and invalidates token if attempts = 5
		updatedRecord, err := s.store.UpdateVerifyAttemptTx(ctx, verifyRecord.VerificationID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if updatedRecord.Attempts == 5 {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("too many invalid requests - please go back to get a new code. [rate_limited]")))
			return
		}

		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("the entered code is incorrect. please try again and check for typos. [digits_mismatch]")))
		return
	}

	//Invalidates token & Creates an User if mode is signup.
	if err := s.store.ManifestTokenTx(ctx, verifyRecord); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(verifyRecord.Email, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := verifyUserResponse{
		AccessToken: accessToken,
		Mode:        verifyRecord.Purpose,
		Email:       verifyRecord.Email,
	}

	ctx.JSON(http.StatusOK, response)
}
