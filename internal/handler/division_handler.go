package handler

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/service"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

type DivisionHandler struct {
	divService service.DivisionService
}

func NewDivisionHandler(divService service.DivisionService) *DivisionHandler {
	return &DivisionHandler{divService: divService}
}

func (h *DivisionHandler) authContext(c *gin.Context) dto.AuthContext {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	return dto.AuthContext{UserID: userID.(string), Role: role.(string)}
}

func (h *DivisionHandler) ListByOrganization(c *gin.Context) {
	res, err := h.divService.ListByOrganization(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *DivisionHandler) Get(c *gin.Context) {
	res, err := h.divService.Get(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *DivisionHandler) Create(c *gin.Context) {
	var req dto.CreateDivisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.divService.Create(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Division created successfully")
}

func (h *DivisionHandler) Update(c *gin.Context) {
	var req dto.UpdateDivisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.divService.Update(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Division updated successfully")
}

func (h *DivisionHandler) Delete(c *gin.Context) {
	if err := h.divService.Delete(h.authContext(c), c.Param("id")); err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, nil, "Division deleted successfully")
}
