package api

import (
	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/IamDushu/Float/internal/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our api service.
type Server struct {
	config util.Config
	store  *db.Store
	router *gin.Engine
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", s.createUser)
	router.POST("/nurses", s.createNurse)
	router.POST("/patients", s.createPatient)
	router.POST("/visits", s.createVisit)

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
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()
	return server, nil
}
