package app

import (
	"net/http"
	"encoding/json"
	"log"
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

type AddReagentReq struct {
    Reagent_id uint `json:"reagent_id" binding:"required"`
    Volume  int  `json:"volume" binding:"required,min=1"`
}

type UpdateReagentVolumeReq struct {
    Volume int `json:"volume" binding:"required,min=1"`
}

// AddReagentToDraft godoc
// @Summary Добавить реагент в черновик заявки
// @Tags methanes
// @Security BearerAuth
// @Param id path int true "Methane ID"
// @Param input body AddReagentReq true "Reagent data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/methanes/{id}/reagents [post]
func (a *Application) AddReagentToDraft(c *gin.Context) {
    id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    userIDAny, ok := c.Get("userid")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    userID := userIDAny.(uint)

    var req AddReagentReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json: " + err.Error()})
        return
    }

    if err := a.repo.AddReagentToDraft(uint(id64), userID, req.Reagent_id, req.Volume); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"ok": true})
}

// UpdateReagentVolume godoc
// @Summary Обновить количество реагента в заявке
// @Tags methanes
// @Security BearerAuth
// @Param id path int true "Methane ID"
// @Param reagent_id path int true "Reagent ID"
// @Param input body AddReagentReq true "Volume"
// @Success 200 {object} map[string]interface{}
// @Router /api/methanes/{id}/reagents/{reagent_id} [put]
func (a Application) UpdateReagentVolume(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid methane id"})
		return
	}

	reagentID64, err := strconv.ParseUint(c.Param("reagent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reagent id"})
		return
	}

	userIDAny, ok := c.Get("userid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	var req UpdateReagentVolumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json: " + err.Error()})
		return
	}

	if err := a.repo.UpdateReagentVolume(uint(id64), userID, uint(reagentID64), req.Volume); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// RemoveReagentFromDraft godoc
// @Summary Удалить реагент из заявки
// @Tags methanes
// @Security BearerAuth
// @Param id path int true "Methane ID"
// @Param reagent_id path int true "Reagent ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/methanes/{id}/reagents/{reagent_id} [delete]
func (a *Application) RemoveReagentFromDraft(c *gin.Context) {
    id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid methane id"})
        return
    }

    reagentID64, err := strconv.ParseUint(c.Param("reagent_id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reagent id"})
        return
    }

    userIDAny, ok := c.Get("userid")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    userID := userIDAny.(uint)

    if err := a.repo.RemoveReagentFromDraft(uint(id64), userID, uint(reagentID64)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"ok": true})
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

    // Проверка сериализации
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

func (a Application) GetMethanesPagedHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	items, total, err := a.repo.GetMethanesPaged(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}