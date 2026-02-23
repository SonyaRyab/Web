package handler

import (
	"lab2/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetReagentById(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	reagent, err := h.Repository.GetReagentByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	// Добавляем полный URL для изображения
	reagent.Img = h.Repository.GetImageURL(reagent.Img)

	ctx.HTML(http.StatusOK, "reagent.page.tmpl", reagent)
}

func (h *Handler) GetAllReagents(ctx *gin.Context) {
	var reagents []ds.Reagent
	var err error

	search := ctx.Query("search")
	if search == "" {
		reagents, err = h.Repository.GetAllReagents()
	} else {
		reagents, err = h.Repository.SearchReagentsByName(search)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	// Добавляем полные URL для изображений
	for i := range reagents {
		reagents[i].Img = h.Repository.GetImageURL(reagents[i].Img)
	}

	ctx.HTML(http.StatusOK, "reagents.page.tmpl", gin.H{
		"data":       reagents,
		"cart_count": h.Repository.GetCartCount(),
		"search":     search,
	})
}

func (h *Handler) AddToExperiment(ctx *gin.Context) {
	reagentIDStr := ctx.PostForm("reagent_id")
	quantityStr := ctx.PostForm("quantity")

	reagentID, err := strconv.ParseUint(reagentIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reagent_id"})
		return
	}

	quantity, err := strconv.ParseFloat(quantityStr, 64)
	if err != nil {
		quantity = 1.0 // по умолчанию
	}

	// Получаем или создаём черновик
	creatorID := uint(1) // захардкодили
	expID, err := h.Repository.GetDraftExperimentID(creatorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Добавляем реагент
	err = h.Repository.AddReagentToExperiment(expID, uint(reagentID), quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/reagents")
}

// GetCurrentExperiment - GET /experiment (ORM)
func (h *Handler) GetCurrentExperiment(ctx *gin.Context) {
	creatorID := uint(1)
	expID, err := h.Repository.GetDraftExperimentID(creatorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exp, items, err := h.Repository.GetExperimentWithReagents(expID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Рассчитываем выход метана
	methaneYield, _ := h.Repository.CalculateMethaneYield(expID)

	ctx.HTML(http.StatusOK, "experiment.page.tmpl", gin.H{
		"experiment":    exp,
		"items":         items,
		"methane_yield": methaneYield,
	})
}

func (h *Handler) DeleteExperiment(ctx *gin.Context) {
	// считываем значение из формы, которую мы добавим в наш шаблон
	expIDStr := ctx.PostForm("experiment_id")
	expID, err := strconv.ParseUint(expIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	err = h.Repository.DeleteExperiment(uint(expID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// после вызова сразу произойдет обновление страницы
	ctx.Redirect(http.StatusFound, "/reagents")
}
