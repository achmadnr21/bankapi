package handler

import (
	"net/http"
	"strconv"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/usecase"
	"github.com/achmadnr21/bankapi/internal/utils"
	"github.com/gin-gonic/gin"
)

type AccountTypeHandler struct {
	uc *usecase.AccountTypeUsecase
}

func NewAccountTypeHandler(uc *usecase.AccountTypeUsecase) *AccountTypeHandler {
	return &AccountTypeHandler{
		uc: uc,
	}
}

func (h *AccountTypeHandler) GetAllAccountTypes(c *gin.Context) {
	accountTypes, err := h.uc.GetAllAccountTypes()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("account types granted", accountTypes))
}

func (h *AccountTypeHandler) GetAccountTypeByID(c *gin.Context) {
	parsed_id := c.Param("id")
	// id kemungkinan string seperti "0011". Ubah menjadi int 11. gunakan convert string to int
	id, err := strconv.Atoi(parsed_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid id"))
		return
	}

	accountType, err := h.uc.GetAccountTypeByID(id)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("account type granted", accountType))
}

func (h *AccountTypeHandler) AddAccountType(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	var accountType domain.AccountType
	if err := c.ShouldBindJSON(&accountType); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}
	newAccountType, err := h.uc.AddAccountType(user_id.(string), accountType)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("account type added", newAccountType))
}

func (h *AccountTypeHandler) UpdateAccountType(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	parsed_id := c.Param("id")
	id, err := strconv.Atoi(parsed_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid id"))
		return
	}
	var accountType domain.AccountType
	if err := c.ShouldBindJSON(&accountType); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}
	updatedAccountType, err := h.uc.UpdateAccountType(user_id.(string), id, accountType)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("account type updated", updatedAccountType))
}
