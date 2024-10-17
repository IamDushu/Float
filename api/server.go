package api

import (
	db "github.com/IamDushu/Float/internal/db/sqlc"
	"github.com/IamDushu/Float/internal/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our api service.
type Server struct {
	config  util.Config
	queries db.Querier
	router  *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, queries db.Querier) (*Server, error) {
	server := &Server{
		config:  config,
		queries: queries,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
