package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetReagentsPublic godoc
// @Summary Список реагентов
// @Description Публичный метод чтения данных
// @Tags reagents
// @Produce json
// @Param search query string false "Поиск"
// @Success 200 {array} ds.Reagent
// @Failure 500 {object} map[string]interface{}
// @Router /api/reagents [get]
func (a *Application) GetReagentsPublic(gCtx *gin.Context) {
	search := gCtx.Query("search")

	items, err := a.repo.GetReagents(search)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gCtx.JSON(http.StatusOK, items)
}

// GetReagentByIDPublic godoc
// @Summary Получить реагент по ID
// @Description Публичный метод чтения одного реагента
// @Tags reagents
// @Produce json
// @Param id path int true "Reagent ID"
// @Success 200 {object} ds.Reagent
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/reagents/{id} [get]
func (a *Application) GetReagentByIDPublic(gCtx *gin.Context) {
	idStr := gCtx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid reagent id"})
		return
	}

	item, err := a.repo.GetReagentByID(uint(id64))
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "reagent not found"})
		return
	}

	gCtx.JSON(http.StatusOK, item)
}