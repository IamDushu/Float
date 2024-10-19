package api

import (
	"net/http"
	"time"

	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createPatientRequest struct {
	UserID                uuid.UUID `json:"user_id" binding:"required"`
	DateOfBirth           time.Time `json:"date_of_birth" binding:"required"`
	EmergencyContactName  string    `json:"emergency_contact_name" binding:"required"`
	EmergencyContactPhone string    `json:"emergency_contact_phone" binding:"required"`
	MedicalHistory        string    `json:"medical_history"`
	Allergies             string    `json:"allergies"`
}

type patientResponse struct {
	PatientID             uuid.UUID `json:"patient_id"`
	UserID                uuid.UUID `json:"user_id"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	EmergencyContactName  string    `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	MedicalHistory        string    `json:"medical_history,omitempty"`
	Allergies             string    `json:"allergies,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}

func newPatientResponse(patient db.Patient) patientResponse {
	return patientResponse{
		PatientID:             patient.PatientID,
		UserID:                patient.UserID,
		DateOfBirth:           patient.DateOfBirth,
		EmergencyContactName:  patient.EmergencyContactName,
		EmergencyContactPhone: patient.EmergencyContactPhone,
		MedicalHistory:        patient.MedicalHistory,
		Allergies:             patient.Allergies,
		CreatedAt:             patient.CreatedAt,
	}
}

type createVisitRequest struct {
	PatientID   uuid.UUID `json:"patient_id" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
	Notes       string    `json:"notes"`
}

type visitResponse struct {
	VisitID     uuid.UUID `json:"visit_id"`
	PatientID   uuid.UUID `json:"patient_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

func newVisitResponse(visit db.Visit) visitResponse {
	return visitResponse{
		VisitID:     visit.VisitID,
		PatientID:   visit.PatientID,
		ScheduledAt: visit.ScheduledAt,
		Status:      visit.Status,
		Notes:       visit.Notes,
		CreatedAt:   visit.CreatedAt,
	}
}

func (s *Server) createPatient(ctx *gin.Context) {
	var req createPatientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreatePatientParams{
		PatientID:             uuid.New(),
		UserID:                req.UserID,
		DateOfBirth:           req.DateOfBirth,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
		MedicalHistory:        req.MedicalHistory,
		Allergies:             req.Allergies,
	}

	patient, err := s.store.CreatePatient(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newPatientResponse(patient)
	ctx.JSON(http.StatusOK, response)
}

func (s *Server) createVisit(ctx *gin.Context) {
	var req createVisitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateVisitParams{
		VisitID:     uuid.New(),
		PatientID:   req.PatientID,
		ScheduledAt: req.ScheduledAt,
		Status:      "pending",
		Notes:       req.Notes,
	}

	visit, err := s.store.CreateVisit(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newVisitResponse(visit)
	ctx.JSON(http.StatusOK, response)
}
