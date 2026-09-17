package handler

import (
	"net/http"

	"ukm-hub/internal/dto"
	"ukm-hub/internal/service"
	"ukm-hub/internal/utils"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	regService service.RegistrationService
}

func NewRegistrationHandler(regService service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{regService: regService}
}

func (h *RegistrationHandler) authContext(c *gin.Context) dto.AuthContext {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	return dto.AuthContext{UserID: userID.(string), Role: role.(string)}
}

func (h *RegistrationHandler) ListByEvent(c *gin.Context) {
	page, perPage := service.ParsePagination(c.Query("page"), c.Query("per_page"))
	res, err := h.regService.ListByEvent(
		h.authContext(c),
		c.Param("id"),
		c.Query("status"),
		c.Query("search"),
		page,
		perPage,
	)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *RegistrationHandler) Get(c *gin.Context) {
	res, err := h.regService.Get(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *RegistrationHandler) Accept(c *gin.Context) {
	res, err := h.regService.Accept(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Applicant accepted successfully")
}

func (h *RegistrationHandler) Reject(c *gin.Context) {
	var req dto.RejectRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.regService.Reject(h.authContext(c), c.Param("id"), req.Reason)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Applicant rejected successfully")
}
