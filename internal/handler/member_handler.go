package handler

import (
	"net/http"

	"handler/internal/dto"
	"handler/internal/service"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService service.MemberService
}

func NewMemberHandler(memberService service.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

func (h *MemberHandler) authContext(c *gin.Context) dto.AuthContext {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	return dto.AuthContext{UserID: userID.(string), Role: role.(string)}
}

func (h *MemberHandler) List(c *gin.Context) {
	page, perPage := service.ParsePagination(c.Query("page"), c.Query("per_page"))
	res, err := h.memberService.List(
		h.authContext(c),
		c.Param("id"),
		c.Query("search"),
		c.Query("division_id"),
		page,
		perPage,
	)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *MemberHandler) Get(c *gin.Context) {
	res, err := h.memberService.Get(h.authContext(c), c.Param("id"))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "")
}

func (h *MemberHandler) Convert(c *gin.Context) {
	var req dto.ConvertMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.memberService.Convert(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, res, "Applicant converted to member successfully")
}

func (h *MemberHandler) Update(c *gin.Context) {
	var req dto.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	res, err := h.memberService.Update(h.authContext(c), c.Param("id"), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, res, "Member updated successfully")
}

func (h *MemberHandler) Delete(c *gin.Context) {
	if err := h.memberService.Delete(h.authContext(c), c.Param("id")); err != nil {
		utils.WriteError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, nil, "Member deleted successfully")
}
