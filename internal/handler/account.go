/*
// Account Router : Private
	accounts := apiV.Group("/accounts")
	accounts.Use(middleware.JWTAuthMiddleware)
	{
		accounts.GET("", accountHandler.GetAllAccounts)
		accounts.GET("/:accnumber", accountHandler.GetAccountByAccountNumber)
		accounts.POST("/:branch/:acctype/:nik", accountHandler.AddAccount)
		accounts.PUT("/:accnumber", accountHandler.UpdateAccount)
		// accounts.DELETE("/:id", branchHandler.DeleteAccount)
	}
*/

package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/usecase"
	"github.com/achmadnr21/bankapi/internal/utils"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	uc *usecase.AccountUsecase
}

func NewAccountHandler(accountUsecase *usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{
		uc: accountUsecase,
	}
}

func (h *AccountHandler) GetAllAccounts(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	accounts, err := h.uc.GetAllAccounts(user_id.(string))
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("accounts granted", accounts))
}

func (h *AccountHandler) GetAccountByAccountNumber(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	parsed_accnumber := c.Param("accnumber")
	// jika account number tidak ada, maka return error
	if parsed_accnumber == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseError("account number is required"))
		return
	}
	account, err := h.uc.GetAccountByAccountNumber(user_id.(string), parsed_accnumber)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("account granted", account))
}

func (h *AccountHandler) AddAccount(c *gin.Context) {
	// get the proposer id
	user_id, _ := c.Get("user_id")
	parsed_nik := c.Param("nik") // the account number is a string
	parsed_branch := c.Param("branch")
	parsed_acctype := c.Param("acctype")

	// convert branch and acctype to int
	branch_id, err := strconv.Atoi(parsed_branch)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid branch id"))
		return
	}
	acctype_id, err := strconv.Atoi(parsed_acctype)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid account type id"))
		return
	}
	// jika nik tidak ada, maka return error
	if parsed_nik == "" {
		c.JSON(http.StatusBadRequest, utils.ResponseError("nik is required"))
		return
	}
	var account domain.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}
	// cetak params
	fmt.Printf("Params:\nbranch_id: %d\nacctype_id: %d\nnik: %s\n", branch_id, acctype_id, parsed_nik)
	// Coba cetak account
	fmt.Printf("1. Account Handler:\n%+v\n", account)
	newAccount, err := h.uc.AddAccount(user_id.(string), branch_id, acctype_id, parsed_nik, account)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.ResponseSuccess("account created", newAccount))
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	// get the proposer id
	user_id, _ := c.Get("user_id")
	parsed_accnumber := c.Param("accnumber") // the account number is a string
	var account domain.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}
	// update the account
	account, err := h.uc.UpdateAccount(user_id.(string), parsed_accnumber, account)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("account updated", account))
}
