/*
	transactionFees := apiV.Group("/transaction-fees")
	transactionFees.Use(middleware.JWTAuthMiddleware)
	{
		transactionFees.GET("", transactionFeeHandler.GetAllTransactionFees)
		transactionFees.GET("/:id", transactionFeeHandler.GetTransactionFeeByID)
		transactionFees.POST("", transactionFeeHandler.AddTransactionFee)
		transactionFees.PUT("/:id", transactionFeeHandler.UpdateTransactionFee)
		transactionFees.DELETE("/:id", transactionFeeHandler.DeleteTransactionFee)
	}
*/

package handler

import (
	"net/http"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/usecase"
	"github.com/achmadnr21/bankapi/internal/utils"
	"github.com/gin-gonic/gin"
)

type TransactionFeeHandler struct {
	uc *usecase.TransactionFeeUsecase
}

func NewTransactionFeeHandler(uc *usecase.TransactionFeeUsecase) *TransactionFeeHandler {
	return &TransactionFeeHandler{
		uc: uc,
	}
}
func (h *TransactionFeeHandler) GetAllTransactionFees(c *gin.Context) {
	transactionFees, err := h.uc.GetAllTransactionFees()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, transactionFees)
}

func (h *TransactionFeeHandler) GetTransactionFeeByID(c *gin.Context) {
	id := c.Param("id")
	transactionFee, err := h.uc.GetTransactionFeeByID(id)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, transactionFee)
}
func (h *TransactionFeeHandler) AddTransactionFee(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not found"))
		return
	}
	var transactionFee domain.TransactionFee
	if err := c.ShouldBindJSON(&transactionFee); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}
	transactionFee, err := h.uc.AddTransactionFee(user_id.(string), transactionFee)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, transactionFee)
}
func (h *TransactionFeeHandler) UpdateTransactionFee(c *gin.Context) {
	id := c.Param("id")
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not found"))
		return
	}
	var transactionFee domain.TransactionFee
	if err := c.ShouldBindJSON(&transactionFee); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}
	transactionFee.ID = id
	transactionFee, err := h.uc.UpdateTransactionFee(user_id.(string), transactionFee)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, transactionFee)
}
func (h *TransactionFeeHandler) DeleteTransactionFee(c *gin.Context) {
	id := c.Param("id")
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("User not found"))
		return
	}
	err := h.uc.DeleteTransactionFee(user_id.(string), id)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusAccepted, utils.ResponseSuccess("Transaction fee deleted", nil))
}
