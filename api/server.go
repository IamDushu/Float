package api

import (
	"fmt"

	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/IamDushu/Float/internal/token"
	"github.com/IamDushu/Float/internal/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our api service.
type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/nurses", s.createNurse)
	router.POST("/patients", s.createPatient)
	router.POST("/visits", s.createVisit)

	router.POST("/api/registration/email", s.registerUser)
	router.POST("/api/registration/email/verify", s.verifyUser)
	router.POST("/api/tokens/renew_access", s.renewAccessToken)

	s.router = router
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}
