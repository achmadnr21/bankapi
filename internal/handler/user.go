package handler

import (
	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/usecase"
	"github.com/achmadnr21/bankapi/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uc *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		uc: userUsecase,
	}
}

func (h *UserHandler) Search(c *gin.Context) {
	// Get the user's id from the URL parameter
	id, _ := c.Get("user_id")
	// The payload is "username", "email", "nik". We need to get it
	// from the request body
	var payload domain.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, utils.ResponseError(err.Error()))
		return
	}
	var username = payload.Username
	var email = payload.Email
	var nik = payload.NIK
	// Call the usecase to search for the user
	users, err := h.uc.Search(id.(string), nik, username, email)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("user search done", users))

}

// Add User
func (h *UserHandler) AddUser(c *gin.Context) {
	id, _ := c.Get("user_id")
	var payload domain.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, utils.ResponseError(err.Error()))
		return
	}
	// Call the usecase to add the user
	user, err := h.uc.AddUser(id.(string), payload)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("user added", user))
}

func (h *UserHandler) GetByNIK(c *gin.Context) {
	nik := c.Param("nik")
	id, _ := c.Get("user_id")
	// Call the usecase to get the user by nik
	user, err := h.uc.GetByNIK(id.(string), nik)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(200, utils.ResponseSuccess("user found", user))
}
