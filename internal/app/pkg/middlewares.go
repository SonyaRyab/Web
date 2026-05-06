package app

import (
	"errors"
	"fmt"
	"log"
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
// func (a Application) WithJWTAuth(assignedRoles ...role.Role) gin.HandlerFunc {
// 	return func(gCtx *gin.Context) {
// 		jwtStr := gCtx.GetHeader("Authorization")
// 		if !strings.HasPrefix(jwtStr, jwtPrefix) {
// 			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
// 			return
// 		}

// 		jwtStr = strings.TrimPrefix(jwtStr, jwtPrefix)

// 		err := a.redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
// 		if err == nil {
// 			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
// 			return
// 		}
// 		if !errors.Is(err, redis.Nil) {
// 			gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
// 			return
// 		}

// 		claims := &ds.JWTClaims{}
// 		token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return []byte(a.config.JWT.Token), nil
// 		})
// 		if err != nil || !token.Valid {
// 			gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
// 			return
// 		}

// 		gCtx.Set("userid", claims.UserID)
// 		gCtx.Set("useruuid", claims.UserUUID)
// 		gCtx.Set("userrole", claims.Role)
// 		gCtx.Set("login", claims.Login)

// 		if len(assignedRoles) == 0 {
// 			gCtx.Next()
// 			return
// 		}

// 		for _, assignedRole := range assignedRoles {
// 			if claims.Role == assignedRole {
// 				gCtx.Next()
// 				return
// 			}
// 		}

// 		gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
// 	}
// }

func (a Application) WithJWTAuth(assignedRoles ...role.Role) gin.HandlerFunc {
    return func(gCtx *gin.Context) {
        jwtStr := gCtx.GetHeader("Authorization")
        log.Printf("[JWT] Authorization header: %s", jwtStr) // +

        if !strings.HasPrefix(jwtStr, jwtPrefix) {
            log.Println("[JWT] Missing bearer prefix") // +
            gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
            return
        }

        jwtStr = strings.TrimPrefix(jwtStr, jwtPrefix)
        log.Printf("[JWT] Token (trimmed): %s", jwtStr[:20]+"...") // +

        err := a.redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
        if err == nil {
            log.Println("[JWT] Token in blacklist") // +
            gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
            return
        }
        if !errors.Is(err, redis.Nil) {
            log.Printf("[JWT] Redis error: %v", err) // +
            gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
            return
        }

        claims := &ds.JWTClaims{}
        token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(a.config.JWT.Token), nil
        })

        if err != nil || !token.Valid {
            log.Printf("[JWT] Parse error: %v, valid=%v", err, token.Valid) // +
            gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        log.Printf("[JWT] Parsed claims: userID=%d, role=%s", claims.UserID, claims.Role) // +

        gCtx.Set("userid", claims.UserID)
        gCtx.Set("useruuid", claims.UserUUID)
        gCtx.Set("userrole", claims.Role)
        gCtx.Set("login", claims.Login)

        if len(assignedRoles) == 0 {
            log.Println("[JWT] No role check required, passing through") // +
            gCtx.Next()
            return
        }

        log.Printf("[JWT] Required roles: %v, user role: %s", assignedRoles, claims.Role) // +

        for _, assignedRole := range assignedRoles {
            if claims.Role == assignedRole {
                log.Println("[JWT] Role matched, passing through") // +
                gCtx.Next()
                return
            }
        }

        log.Println("[JWT] Role not matched, forbidden") // +
        gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
    }
}