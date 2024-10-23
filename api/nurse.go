package api

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/IamDushu/Float/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createNurseRequest struct {
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required,min=6"`
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	PhoneNumber       string `json:"phone_number" binding:"required"`
	LicenseNumber     string `json:"license_number"  binding:"required"`
	Specialization    string `json:"specialization"  binding:"required"`
	YearsOfExperience int32  `json:"years_of_experience"  binding:"required"`
	ZipCode           string `json:"zip_code"  binding:"required,max=5"`
}

type nurseResponse struct {
	NurseID           uuid.UUID `json:"nurse_id"`
	UserID            uuid.UUID `json:"user_id"`
	LicenseNumber     string    `json:"license_number"`
	Specialization    string    `json:"specialization"`
	YearsOfExperience int32     `json:"years_of_experience"`
	ZipCode           string    `json:"zip_code"`
	CreatedAt         time.Time `json:"created_at"`
}

type nurseAccountResponse struct {
	User  userResponse
	Nurse nurseResponse
}

func newNurseResponse(nurse db.CreateNurseAccountResult) nurseAccountResponse {
	return nurseAccountResponse{
		User: userResponse{
			UserID:      nurse.User.UserID,
			Email:       nurse.User.Email,
			FirstName:   nurse.User.FirstName,
			LastName:    nurse.User.LastName,
			PhoneNumber: nurse.User.PhoneNumber,
			CreatedAt:   nurse.User.CreatedAt,
		},
		Nurse: nurseResponse{
			NurseID:           nurse.Nurse.NurseID,
			UserID:            nurse.User.UserID,
			LicenseNumber:     nurse.Nurse.LicenseNumber,
			Specialization:    nurse.Nurse.Specialization,
			YearsOfExperience: nurse.Nurse.YearsOfExperience,
			ZipCode:           nurse.Nurse.ZipCode,
			CreatedAt:         nurse.Nurse.CreatedAt,
		},
	}
}

func (s *Server) createNurse(ctx *gin.Context) {
	var req createNurseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phoneDetails, err := util.VerifyPhone(req.PhoneNumber, s.config.TwillioAccountSID, s.config.TwillioAuthToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !phoneDetails.Valid {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid phone number: %v", phoneDetails.ValidationErrors)))
		return
	}

	arg := db.CreateNurseAccountParams{
		Email:             req.Email,
		PasswordHash:      req.Password,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		PhoneNumber:       phoneDetails.PhoneNumber,
		LicenseNumber:     req.LicenseNumber,
		Specialization:    req.Specialization,
		YearsOfExperience: req.YearsOfExperience,
		ZipCode:           req.ZipCode,
	}

	nurse, err := s.store.CreateNurseAccountTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newNurseResponse(nurse)
	ctx.JSON(http.StatusOK, response)
}
