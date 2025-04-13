package handler

import (
	"net/http"
	"strconv"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/usecase"
	"github.com/achmadnr21/bankapi/internal/utils"
	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	uc *usecase.BranchUsecase
}

func NewBranchHandler(uc *usecase.BranchUsecase) *BranchHandler {
	return &BranchHandler{
		uc: uc,
	}
}

func (h *BranchHandler) GetAllBranches(c *gin.Context) {
	branches, err := h.uc.GetAllBranches()
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("branches granted", branches))
}

func (h *BranchHandler) GetBranchByID(c *gin.Context) {
	parsed_id := c.Param("id")
	id, err := strconv.Atoi(parsed_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid id"))
		return
	}
	branch, err := h.uc.GetBranchByID(id)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("branch granted", branch))
}

func (h *BranchHandler) AddBranch(c *gin.Context) {
	// get the proposer id
	id, _ := c.Get("user_id")
	var branch domain.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}

	newBranch, err := h.uc.AddBranch(id.(string), branch)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.ResponseSuccess("branch added", newBranch))
}

func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	// get the proposer id
	user_id, _ := c.Get("user_id")
	parsed_id := c.Param("id")
	id, err := strconv.Atoi(parsed_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid id"))
		return
	}

	var branch domain.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseError("Invalid request"))
		return
	}

	updatedBranch, err := h.uc.UpdateBranch(user_id.(string), id, branch)
	if err != nil {
		c.JSON(utils.GetHTTPErrorCode(err), utils.ResponseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.ResponseSuccess("branch updated", updatedBranch))
}
