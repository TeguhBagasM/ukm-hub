package handler

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/service"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	formService service.FormService
}

func NewFormHandler(formService service.FormService) *FormHandler {
	return &FormHandler{formService: formService}
}

func (h *FormHandler) authContext(c *gin.Context) dto.AuthContext {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	return dto.AuthContext{UserID: userID.(string), Role: role.(string)}
}

func (h *FormHandler) GetByEvent(c *gin.Context) {
	res, err := h.formService.GetByEvent(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *FormHandler) CreateForm(c *gin.Context) {
	var req dto.CreateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.formService.CreateForm(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Form created successfully")
}

func (h *FormHandler) UpdateForm(c *gin.Context) {
	var req dto.UpdateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.formService.UpdateForm(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Form updated successfully")
}

func (h *FormHandler) Publish(c *gin.Context) {
	res, err := h.formService.Publish(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Form published successfully")
}

func (h *FormHandler) Unpublish(c *gin.Context) {
	res, err := h.formService.Unpublish(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Form unpublished successfully")
}

func (h *FormHandler) AddField(c *gin.Context) {
	var req dto.FormFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.formService.AddField(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Field added successfully")
}

func (h *FormHandler) UpdateField(c *gin.Context) {
	var req dto.FormFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.formService.UpdateField(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Field updated successfully")
}

func (h *FormHandler) DeleteField(c *gin.Context) {
	if err := h.formService.DeleteField(h.authContext(c), c.Param("id")); err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, nil, "Field deleted successfully")
}

func (h *FormHandler) DuplicateField(c *gin.Context) {
	res, err := h.formService.DuplicateField(h.authContext(c), c.Param("id"), c.Param("fieldId"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Field duplicated successfully")
}

func (h *FormHandler) ReorderFields(c *gin.Context) {
	var items []dto.ReorderFieldRequest
	if err := c.ShouldBindJSON(&items); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.formService.ReorderFields(h.authContext(c), c.Param("id"), items); err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, nil, "Fields reordered successfully")
}
