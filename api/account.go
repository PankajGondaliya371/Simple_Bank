package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqerror, ok := err.(*pq.Error); ok {
			switch pqerror.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}

// Get request validation
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, request.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)
}

// List account request validation
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var request listAccountRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.LISTACCOUNTSParams{
		Owner:  authPayload.Username,
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	accounts, err := server.store.LISTACCOUNTS(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
