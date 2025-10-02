package api

import (
	"database/sql"
	"net/http"

	db "github.com/ShubhKanodia/GoBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// what is ctx *gin.Context? -> it is a context object that contains information about the HTTP request and response.
// It is used to access request data, set response data, and handle errors.
// It is passed to the handler functions to provide access to the request and response objects.
//provides convenience methods for working with HTTP requests and reading input params.

// CreateAccountRequest contains the request parameters for creating a new account.
// Server serves HTTP requests for our banking service.

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"` // uri:"id" means that the id is extracted from the URL path
	// binding:"required,min=1" means that the id is required and must be greater than or equal to 1
}

type listAccountsRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"` //form tage is used for query parameters
	// PageID is the page number for pagination, starting from 1
	//pageid and pagesize are query parameters that are used for pagination, whereas id is a uri parameter
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	//Bind validation rules to the request body
	// If the request body is not valid, it will return a 400 Bad Request error
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		// ShouldBindJSON binds the request body to the struct and validates it
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx.Request.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// Check if the error is a PostgreSQL error
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				// Handle specific PostgreSQL error codes
				ctx.JSON(http.StatusConflict, gin.H{"error": pqErr.Message}) // 409 Conflict for unique constraint violations
			default:
				// For other PostgreSQL errors, return a generic internal server error
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": pqErr.Message})

			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	// Bind the request parameters to the struct
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id := req.ID
	account, err := server.store.GetAccount(ctx.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) ListAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	// Bind the request parameters to the struct
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// List accounts with pagination

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize, // Offset is calculated based on the page number and page size
	}
	accounts, err := server.store.ListAccounts(ctx.Request.Context(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
