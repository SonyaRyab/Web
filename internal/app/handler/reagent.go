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

	// Добавляем полный URL для изображения и видео
	reagent.Img = h.Repository.GetImageURL(reagent.Img)
	reagent.Video = h.Repository.GetImageURL(reagent.Video)

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
		reagents[i].Video = h.Repository.GetImageURL(reagents[i].Video)
	}

	ctx.HTML(http.StatusOK, "reagents.page.tmpl", gin.H{
		"data":       reagents,
		"cart_count": h.Repository.GetCartCount(),
		"search":     search,
	})
}

func (h *Handler) AddToMethane(ctx *gin.Context) {
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
	adminID := uint(1) // захардкодили
	expID, err := h.Repository.GetDraftMethaneID(adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Добавляем реагент
	err = h.Repository.AddReagentToMethane(expID, uint(reagentID), quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/reagents")
}

// GetCurrentExperiment - GET /experiment (ORM)
func (h *Handler) GetCurrentMethane(ctx *gin.Context) {
	adminID := uint(1)
	methaneID, err := h.Repository.GetDraftMethaneID(adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	methane, items, err := h.Repository.GetMethaneWithReagents(methaneID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Рассчитываем выход метана
	methaneYield, _ := h.Repository.CalculateMethaneYield(methaneID)

	ctx.HTML(http.StatusOK, "methane.page.tmpl", gin.H{
		"methane":       methane,
		"items":         items,
		"methane_yield": methaneYield,
	})
}

func (h *Handler) DeleteMethane(ctx *gin.Context) {
	// считываем значение из формы, которую мы добавим в наш шаблон
	methaneIDStr := ctx.PostForm("methane_id")
	methaneID, err := strconv.ParseUint(methaneIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	err = h.Repository.DeleteMethane(uint(methaneID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// после вызова сразу произойдет обновление страницы
	ctx.Redirect(http.StatusFound, "/reagents")
}
