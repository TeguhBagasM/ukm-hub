package handler

import (
	"fmt"
	"net/http"

	"handler/internal/dto"
	"handler/internal/service"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) authContext(c *gin.Context) dto.AuthContext {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	return dto.AuthContext{UserID: userID.(string), Role: role.(string)}
}

func (h *DashboardHandler) Get(c *gin.Context) {
	res, err := h.dashboardService.GetDashboard(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *DashboardHandler) ExportRegistrations(c *gin.Context) {
	filename, content, err := h.dashboardService.ExportRegistrations(h.authContext(c), c.Param("id"), c.Query("event_id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	writeCSV(c, filename, content)
}

func (h *DashboardHandler) ExportMembers(c *gin.Context) {
	filename, content, err := h.dashboardService.ExportMembers(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	writeCSV(c, filename, content)
}

func writeCSV(c *gin.Context, filename string, content []byte) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", content)
}
