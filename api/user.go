package api

import (
	"net/http"
	"time"

	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/IamDushu/Float/internal/util"
	normalizer "github.com/dimuska139/go-email-normalizer/v3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"max=10"`
}

type userResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	var phone_number *string
	if user.PhoneNumber.Valid {
		phone_number = &user.PhoneNumber.String
	}
	return userResponse{
		UserID:      user.UserID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: phone_number,
		CreatedAt:   user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	n := normalizer.NewNormalizer()
	req.Email = n.Normalize(req.Email)

	arg := db.CreateUserParams{
		UserID:       uuid.New(),
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: req.Password,
		PhoneNumber:  util.ToNullString(req.PhoneNumber),
	}

	user, err := s.queries.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}
