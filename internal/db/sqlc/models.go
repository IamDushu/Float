// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	VerificationID uuid.UUID `json:"verification_id"`
	Email          string    `json:"email"`
	Token          string    `json:"token"`
	HashedOtp      string    `json:"hashed_otp"`
	Purpose        string    `json:"purpose"`
	Attempts       int32     `json:"attempts"`
	ExpiresAt      time.Time `json:"expires_at"`
	Valid          bool      `json:"valid"`
	CreatedAt      time.Time `json:"created_at"`
}

type Nurse struct {
	NurseID           uuid.UUID `json:"nurse_id"`
	UserID            uuid.UUID `json:"user_id"`
	LicenseNumber     string    `json:"license_number"`
	Specialization    string    `json:"specialization"`
	YearsOfExperience int32     `json:"years_of_experience"`
	ZipCode           string    `json:"zip_code"`
	CreatedAt         time.Time `json:"created_at"`
}

type NurseAvailability struct {
	AvailabilityID uuid.UUID `json:"availability_id"`
	NurseID        uuid.UUID `json:"nurse_id"`
	DayOfWeek      string    `json:"day_of_week"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	CreatedAt      time.Time `json:"created_at"`
}

type Patient struct {
	PatientID             uuid.UUID `json:"patient_id"`
	UserID                uuid.UUID `json:"user_id"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	EmergencyContactName  string    `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	MedicalHistory        string    `json:"medical_history"`
	Allergies             string    `json:"allergies"`
	CreatedAt             time.Time `json:"created_at"`
}

type Session struct {
	SessionID    uuid.UUID `json:"session_id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type Visit struct {
	VisitID     uuid.UUID     `json:"visit_id"`
	NurseID     uuid.NullUUID `json:"nurse_id"`
	PatientID   uuid.UUID     `json:"patient_id"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	CompletedAt time.Time     `json:"completed_at"`
	Status      string        `json:"status"`
	Notes       string        `json:"notes"`
	CreatedAt   time.Time     `json:"created_at"`
}
