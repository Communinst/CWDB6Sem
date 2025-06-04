package handler

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		tempId := c.MustGet("user-id")
		if tempId == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			c.Abort()
			return
		}

		userId, err := strconv.Atoi(tempId.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		user, err := h.service.AuthServiceInterface.GetUser(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		defaultAdmin, _ := strconv.ParseInt(os.Getenv("DEFAULT_ADMIN_ROLE_ID"), 10, 64)

		if user.RoleId != int(defaultAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admin access required"})
			c.Abort()
			return
		}

		c.Set("admin", true)
		c.Next()
	}
}
