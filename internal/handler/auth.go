package handler

import (
	"net/http"
	"strings"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/usecase"
	"github.com/achmadnr21/bankapi/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		uc: authUsecase,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError(err.Error()))
		return
	}

	token, refreshToken, err := h.uc.Login(user.Username, user.Password)
	if err != nil {
		var errCode int = utils.GetHTTPErrorCode(err)
		c.JSON(errCode, utils.ResponseError(err.Error()))
		return
	}

	var tokenResponse = struct {
		Token   string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}{
		Token:   token,
		Refresh: refreshToken,
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Login success", tokenResponse))

}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	tokenString := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, utils.ResponseError("Token not found"))
		c.Abort()
		return
	}
	// proses validasi di usecase
	token, refreshToken, err := h.uc.RefreshToken(tokenString)
	if err != nil {
		var errCode int = utils.GetHTTPErrorCode(err)
		c.JSON(errCode, utils.ResponseError(err.Error()))
		return
	}
	var tokenResponse = struct {
		Token   string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}{
		Token:   token,
		Refresh: refreshToken,
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("Token refreshed", tokenResponse))
}
