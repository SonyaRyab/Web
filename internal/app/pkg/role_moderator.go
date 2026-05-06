package app

import (
    "net/http"

    "lab4/internal/app/role"

    "github.com/gin-gonic/gin"
)

func (a Application) RequireModerator() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		roleAny, exists := gCtx.Get("userrole")
		if !exists {
			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		currentRole := roleAny.(role.Role)
		if currentRole != role.Moderator && currentRole != role.Admin {
			gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		gCtx.Next()
	}
}
