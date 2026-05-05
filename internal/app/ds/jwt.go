package ds

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"lab4/internal/app/role"
)

type JWTClaims struct {
	jwt.StandardClaims           // все что точно необходимо по RFC
	UserID   uint      `json:"user_id"`
	UserUUID           uuid.UUID `json:"user_uuid"` // наши данные - uuid этого пользователя в базе данных
	Login    		   string    `json:"login"`
	Role               role.Role `json:"role"`
	Scopes             []string  `json:"scopes" json:"scopes"` // список доступов в нашей системе
}