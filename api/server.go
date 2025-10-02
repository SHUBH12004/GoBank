package api

import (
	db "github.com/ShubhKanodia/GoBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	// Store is the database store that provides access to the database
	store db.Store
	// Store is the database store that provides access to the database

	router *gin.Engine
	// Router is the HTTP router that handles incoming requests
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	//register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	router.POST("/accounts", server.createAccount)   // last one shoukd be the handler func like (middleware1, middleware2..., handler)
	router.GET("/accounts/:id", server.getAccount)   // this is the endpoint for getting an account by id
	router.GET("/accounts", server.ListAccounts)     // this is the endpoint for listing accounts with pagination
	router.POST("/transfers", server.createTransfer) // this is the endpoint for creating a transfer
	server.router = router
	return server
}
func (s *Server) Start(address string) error {
	// Start the HTTP server on the given address
	return s.router.Run(address)
}

// errorResponse creates a JSON response for errors
// gin .H
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
