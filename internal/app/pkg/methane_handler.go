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

// GetMethaneByID godoc
// @Summary Получить заявку по ID
// @Tags methanes
// @Security BearerAuth
// @Param id path int true "Methane ID"
// @Success 200 {object} ds.Methane
// @Failure 404 {object} map[string]interface{}
// @Router /api/methanes/{id} [get]
Да, тут уже почти наверняка не роут и не фронт, а именно ошибка сериализации JSON на бэкенде: Gin может отдать 200 OK с пустым body, если c.JSON(...) не смог закодировать структуру в JSON. Это типичное поведение при “плохом” значении внутри ответа, например NaN, Inf или проблемном вложенном поле.

Что это значит у тебя
Раз GET /api/methanes/8 возвращает 200, но response: null, значит a.repo.GetMethaneByID(...) заявку находит, однако ломается именно на c.JSON(http.StatusOK, methane).

Скорее всего проблема в одном из полей структуры ds.Methane или вложенных структур:

Temperature / MethaneYield могут содержать NaN или Inf;

в Reagents или Reagents.Reagent может быть поле, которое не сериализуется;

либо во вложенном User/другой модели есть значение, которое ломает marshal.

Сделай точную диагностику
Прямо в GetMethaneByID добавь проверку через стандартный encoding/json, до c.JSON(...). Это сразу покажет реальную ошибку.

Замени handler на такой:

go
package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	_, err = json.Marshal(methane)
	if err != nil {
		log.Println("GetMethaneByID marshal error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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