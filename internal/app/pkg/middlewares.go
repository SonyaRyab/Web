package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"lab4/internal/app/ds"
	"lab4/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

const jwtPrefix = "Bearer "

// WithJWTAuth - middleware для проверки JWT и ролей
func (a *Application) WithJWTAuth(assignedRoles ...role.Role) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		
		jwtStr = strings.TrimPrefix(jwtStr, jwtPrefix)

		// Проверка blacklist в Redis
		err := a.redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil {
			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}
		if !errors.Is(err, redis.Nil) {
			gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
			return
		}

		// Парсинг токена
		claims := &ds.JWTClaims{}
		token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.config.JWT.Token), nil
		})
		
		if err != nil || !token.Valid {
			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Устанавливаем данные пользователя в контекст
		gCtx.Set("userid", claims.UserID)
		gCtx.Set("useruuid", claims.UserUUID)
		gCtx.Set("userrole", claims.Role)
		gCtx.Set("login", claims.Login)

		// Проверка ролей
		if len(assignedRoles) == 0 {
			gCtx.Next()
			return
		}

		for _, assignedRole := range assignedRoles {
			if claims.Role == assignedRole {
				gCtx.Next()
				return
			}
		}

		gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}