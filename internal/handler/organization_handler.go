package handler

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/service"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgService service.OrganizationService
}

func NewOrganizationHandler(orgService service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

func (h *OrganizationHandler) authContext(c *gin.Context) dto.AuthContext {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	return dto.AuthContext{UserID: userID.(string), Role: role.(string)}
}

func (h *OrganizationHandler) List(c *gin.Context) {
	res, err := h.orgService.List(h.authContext(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *OrganizationHandler) Get(c *gin.Context) {
	res, err := h.orgService.Get(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.orgService.Create(h.authContext(c), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Organization created successfully")
}

func (h *OrganizationHandler) Update(c *gin.Context) {
	var req dto.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.orgService.Update(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Organization updated successfully")
}

func (h *OrganizationHandler) Delete(c *gin.Context) {
	if err := h.orgService.Delete(h.authContext(c), c.Param("id")); err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, nil, "Organization deleted successfully")
}
