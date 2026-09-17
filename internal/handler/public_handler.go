package handler

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/service"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	publicService service.PublicService
}

func NewPublicHandler(publicService service.PublicService) *PublicHandler {
	return &PublicHandler{publicService: publicService}
}

func (h *PublicHandler) GetEvent(c *gin.Context) {
	res, err := h.publicService.GetEvent(c.Param("slug"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *PublicHandler) GetForm(c *gin.Context) {
	res, err := h.publicService.GetForm(c.Param("slug"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *PublicHandler) Register(c *gin.Context) {
	var req dto.PublicRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.publicService.Register(c.Param("slug"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Registration submitted successfully")
}
