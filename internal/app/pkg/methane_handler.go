package app

import (
	"net/http"
	"strconv"

	"lab4/internal/app/role"

	"github.com/gin-gonic/gin"
)

type FormMethaneReq struct {
	Name         string  `json:"name"`
	Temperature  float64 `json:"temperature"`
	MethaneYield float64 `json:"methane_yield"`
}

type CompleteMethaneReq struct {
	Status string `json:"status"`
}

func (a *Application) GetMethanes(gCtx *gin.Context) {
	userIDAny, ok := gCtx.Get("userid")
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRoleAny, ok := gCtx.Get("userrole")
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDAny.(uint)
	userRole := userRoleAny.(role.Role)

	items, err := a.repo.GetMethanesForUser(userID, userRole)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, items)
}

func (a *Application) GetMethaneByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid methane id"})
		return
	}

	methane, err := a.repo.GetMethaneByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "methane not found"})
		return
	}

	c.JSON(http.StatusOK, methane)
}

func (a *Application) CreateDraftMethane(gCtx *gin.Context) {
	userIDAny, ok := gCtx.Get("userid")
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDAny.(uint)

	item, err := a.repo.CreateDraftMethane(userID)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, item)
}

func (a *Application) GetDraftMethane(gCtx *gin.Context) {
	userIDAny, ok := gCtx.Get("userid")
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDAny.(uint)

	item, err := a.repo.GetDraftMethane(userID)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, item)
}

func (a *Application) FormMethane(gCtx *gin.Context) {
	id64, err := strconv.ParseUint(gCtx.Param("id"), 10, 64)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userIDAny, ok := gCtx.Get("userid")
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	var req FormMethaneReq
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	updates := map[string]interface{}{
		"name":          req.Name,
		"temperature":   req.Temperature,
		"methane_yield": req.MethaneYield,
	}

	if err := a.repo.FormMethaneByOwner(uint(id64), userID, updates); err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (a *Application) CompleteMethane(gCtx *gin.Context) {
	id64, err := strconv.ParseUint(gCtx.Param("id"), 10, 64)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userIDAny, ok := gCtx.Get("userid")
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUserID := userIDAny.(uint)

	var req CompleteMethaneReq
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := a.repo.CompleteMethane(uint(id64), currentUserID, req.Status); err != nil {
		if err.Error() == "forbidden" {
			gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, gin.H{"ok": true})
}