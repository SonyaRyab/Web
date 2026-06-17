package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parsePositiveInt(value string, def int) int {
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 {
		return def
	}
	return n
}

// GetReagentsPublic godoc
// @Summary Список реагентов
// @Description Публичный метод чтения данных
// @Tags reagents
// @Produce json
// @Param search query string false "Поиск"
// @Success 200 {array} ds.Reagent
// @Failure 500 {object} map[string]interface{}
// @Router /api/reagents [get]
func (a Application) GetReagentsPublic(c *gin.Context) {
	search := c.Query("search")
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	limit := parsePositiveInt(c.DefaultQuery("limit", "24"), 24)

	if limit > 100 {
		limit = 100
	}

	items, total, err := a.repo.GetReagentsPaged(search, page, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := 1
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      items,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
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
func (a Application) GetReagentByIDPublic(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid reagent id"})
		return
	}

	item, err := a.repo.GetReagentByID(uint(id64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "reagent not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}