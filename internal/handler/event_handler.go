package handler

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/service"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) authContext(c *gin.Context) dto.AuthContext {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	return dto.AuthContext{UserID: userID.(string), Role: role.(string)}
}

func (h *EventHandler) ListByOrganization(c *gin.Context) {
	res, err := h.eventService.ListByOrganization(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *EventHandler) Get(c *gin.Context) {
	res, err := h.eventService.Get(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.eventService.Create(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Event created successfully")
}

func (h *EventHandler) Update(c *gin.Context) {
	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.eventService.Update(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Event updated successfully")
}

func (h *EventHandler) Delete(c *gin.Context) {
	if err := h.eventService.Delete(h.authContext(c), c.Param("id")); err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, nil, "Event deleted successfully")
}

func (h *EventHandler) Publish(c *gin.Context) {
	res, err := h.eventService.Publish(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Event published successfully")
}

func (h *EventHandler) Close(c *gin.Context) {
	res, err := h.eventService.Close(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Event closed successfully")
}

func (h *EventHandler) Archive(c *gin.Context) {
	res, err := h.eventService.Archive(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Event archived successfully")
}
