package middleware

import (
	"net/http"

	"ukm-hub/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) func(c *gin.Context) {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		roleStr, _ := role.(string)
		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}
		utils.ErrorResponse(c, http.StatusForbidden, "forbidden: insufficient role")
		c.Abort()
	}
}
