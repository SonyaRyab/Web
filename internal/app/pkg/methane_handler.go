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

// GetMethanes godoc
// @Summary Список заявок
// @Description Для исследователя возвращает только его заявки, для модератора и администратора — все
// @Tags methanes
// @Produce json
// @Success 200 {array} ds.Methane
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/methanes [get]
func (a *Application) GetMethanes(gCtx *gin.Context) {
	userIDAny, ok := gCtx.Get("user_id")
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRoleAny, ok := gCtx.Get("user_role")
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

// CreateDraftMethane godoc
// @Summary Создать черновик заявки
// @Description Создаёт новую заявку и назначает текущего пользователя автором
// @Tags methanes
// @Produce json
// @Success 200 {object} ds.Methane
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/methanes/draft [post]
func (a *Application) CreateDraftMethane(gCtx *gin.Context) {
	userIDAny, ok := gCtx.Get("user_id")
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

// GetDraftMethane godoc
// @Summary Получить черновик текущего пользователя
// @Tags methanes
// @Produce json
// @Success 200 {object} ds.Methane
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/methanes/draft [get]
func (a *Application) GetDraftMethane(gCtx *gin.Context) {
	userIDAny, ok := gCtx.Get("user_id")
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

// FormMethane godoc
// @Summary Сформировать заявку
// @Description Только владелец заявки может перевести её из черновика в статус "сформирована"
// @Tags methanes
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param input body FormMethaneReq true "Поля заявки для обновления"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/methanes/{id}/form [put]
func (a *Application) FormMethane(gCtx *gin.Context) {
	id64, err := strconv.ParseUint(gCtx.Param("id"), 10, 64)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userIDAny, ok := gCtx.Get("user_id")
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

// CompleteMethane godoc
// @Summary Завершить или отклонить заявку
// @Description Только модератор или администратор может завершить сформированную заявку
// @Tags methanes
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param input body CompleteMethaneReq true "Статус завершения"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/methanes/{id}/complete [put]
func (a *Application) CompleteMethane(gCtx *gin.Context) {
	id64, err := strconv.ParseUint(gCtx.Param("id"), 10, 64)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userIDAny, ok := gCtx.Get("user_id")
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
