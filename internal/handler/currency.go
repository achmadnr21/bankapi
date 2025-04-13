/*
// Add Currency Router : Private
currencies := apiV.Group("/currencies")
currencies.Use(middleware.JWTAuthMiddleware)
{
	currencies.GET("", currencyHandler.GetAllCurrencies)
	currencies.GET("/:id", currencyHandler.GetCurrencyByID)
	currencies.POST("", currencyHandler.AddCurrency)
	currencies.PUT("/:id", currencyHandler.UpdateCurrency)
	// currencies.DELETE("/:id", branchHandler.DeleteCurrency)
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

type CurrencyHandler struct {
	uc *usecase.CurrencyUsecase
}

func NewCurrencyHandler(currencyUsecase *usecase.CurrencyUsecase) *CurrencyHandler {
	return &CurrencyHandler{uc: currencyUsecase}
}

func (h *CurrencyHandler) GetAllCurrencies(c *gin.Context) {
	currencies, err := h.uc.GetAllCurrencies()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("currencies retrieved", currencies))
}

func (h *CurrencyHandler) GetCurrencyByID(c *gin.Context) {
	id := c.Param("id")
	currency, err := h.uc.GetCurrencyByID(id)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("currency retrieved", currency))
}

func (h *CurrencyHandler) AddCurrency(c *gin.Context) {
	// Get the user's id from the URL parameter
	user_id, _ := c.Get("user_id")
	var payload domain.Currency
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, utils.ResponseError(err.Error()))
		return
	}
	currency, err := h.uc.AddCurrency(user_id.(string), payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	// Return the created currency
	c.JSON(http.StatusCreated, utils.ResponseSuccess("currency added", currency))
}

func (h *CurrencyHandler) UpdateCurrency(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	id := c.Param("id")
	var payload domain.Currency
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, utils.ResponseError(err.Error()))
		return
	}
	currency, err := h.uc.UpdateCurrency(user_id.(string), id, payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("currency updated", currency))
}
