package app

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lab4/internal/app/ds"
	"lab4/internal/app/role"
)

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	Login       string `json:"login"`
	Username    string `json:"username"`
	ExpiresIn   int64  `json:"expiresin"`
	AccessToken string `json:"accesstoken"`
	TokenType   string `json:"tokentype"`
}

type registerReq struct {
	Login    string `json:"login"`
	Username string `json:"name"`
	Pass     string `json:"pass"`
}

type registerResp struct {
	Ok bool `json:"ok"`
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Register godoc
// @Summary
// @Description
// @Tags auth
// @Accept json
// @Produce json
// @Param input body registerReq true "register request"
// @Success 200 {object} registerResp
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (a Application) Register(gCtx *gin.Context) {
	var req registerReq
	if err := json.NewDecoder(gCtx.Request.Body).Decode(&req); err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if req.Login == "" || req.Username == "" || req.Pass == "" {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "all fields required"})
		return
	}

	_, err := a.repo.GetUserByLogin(req.Login)
	if err == nil {
		gCtx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "login already exists"})
		return
	}

	passHash, err := hashPassword(req.Pass)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot hash password"})
		return
	}

	err = a.repo.Register(&ds.User{
	UUID:     uuid.New(),
	Login:    req.Login,
	Username: req.Username,
	Email:    req.Login + "@test.ru",
	Role:     role.Researcher,
	PassHash: passHash,
	})
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, registerResp{Ok: true})
}

// Login godoc
// @Summary JWT
// @Description JWT, accesstoken
// @Tags auth
// @Accept json
// @Produce json
// @Param input body loginReq true "login request"
// @Success 200 {object} loginResp
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (a Application) Login(gCtx *gin.Context) {
	var req loginReq
	if err := json.NewDecoder(gCtx.Request.Body).Decode(&req); err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	user, err := a.repo.GetUserByLogin(req.Login)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !checkPassword(user.PassHash, req.Password) {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	now := time.Now()
	expiresAt := now.Add(a.config.JWT.ExpiresIn).Unix()

	token := jwt.NewWithClaims(
		jwt.GetSigningMethod(a.config.JWT.SigningMethod),
		ds.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiresAt,
				IssuedAt:  now.Unix(),
				Issuer:    "lab4-backend",
			},
			UserID:   user.ID,
			UserUUID: user.UUID,
			Login:    user.Login,
			Role:     user.Role,
		},
	)

	tokenString, err := token.SignedString([]byte(a.config.JWT.Token))
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot create token"})
		return
	}

	gCtx.JSON(http.StatusOK, loginResp{
		Login:       user.Login,
		Username:    user.Username,
		ExpiresIn:   int64(a.config.JWT.ExpiresIn.Seconds()),
		AccessToken: tokenString,
		TokenType:   "Bearer",
	})
}

// Logout godoc
// @Summary Blacklist JWT
// @Description JWT blacklist Redis
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/logout [post]
func (a Application) Logout(gCtx *gin.Context) {
	jwtStr := gCtx.GetHeader("Authorization")
	const prefix = "Bearer "

	if !strings.HasPrefix(jwtStr, prefix) {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid auth header"})
		return
	}

	jwtStr = strings.TrimPrefix(jwtStr, prefix)

	claims := &ds.JWTClaims{}
	_, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.JWT.Token), nil
	})
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	expiration := time.Unix(claims.ExpiresAt, 0)
	ttl := time.Until(expiration)
	if ttl <= 0 {
		gCtx.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	if err := a.redis.WriteJWTToBlacklist(gCtx.Request.Context(), jwtStr, ttl); err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot blacklist token"})
		return
	}

	gCtx.JSON(http.StatusOK, gin.H{"ok": true})
}