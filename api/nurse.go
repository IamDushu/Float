package api

import (
	"net/http"
	"time"

	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createNurseRequest struct {
	UserID            uuid.UUID `json:"user_id" binding:"required"`
	LicenseNumber     string    `json:"license_number"  binding:"required"`
	Specialization    string    `json:"specialization"  binding:"required"`
	YearsOfExperience int32     `json:"years_of_experience"  binding:"required"`
	ZipCode           string    `json:"zip_code"  binding:"required,max=5"`
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

func newNurseResponse(nurse db.Nurse) nurseResponse {
	return nurseResponse{
		NurseID:           nurse.NurseID,
		UserID:            nurse.UserID,
		LicenseNumber:     nurse.LicenseNumber,
		Specialization:    nurse.Specialization,
		YearsOfExperience: nurse.YearsOfExperience,
		ZipCode:           nurse.ZipCode,
		CreatedAt:         nurse.CreatedAt,
	}
}

func (s *Server) createNurse(ctx *gin.Context) {
	var req createNurseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateNurseParams{
		NurseID:           uuid.New(),
		UserID:            req.UserID,
		LicenseNumber:     req.LicenseNumber,
		Specialization:    req.Specialization,
		YearsOfExperience: req.YearsOfExperience,
		ZipCode:           req.ZipCode,
	}

	nurse, err := s.store.CreateNurse(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newNurseResponse(nurse)
	ctx.JSON(http.StatusOK, response)
}
