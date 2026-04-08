package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddReagentToMethaneAPI(ctx *gin.Context) {
	methaneID, _ := strconv.Atoi(ctx.Param("id"))

	var input struct {
		ReagentID uint    `json:"reagent_id" binding:"required"`
		Quantity  float64 `json:"quantity" binding:"required"`
		OrderNum  int     `json:"order_num"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверяем, что заявка в статусе черновик
	methane, err := h.Repository.GetMethane(methaneID)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	if methane.Status != "черновик" {
		h.errorHandler(ctx, http.StatusBadRequest,
			fmt.Errorf("можно добавлять реагенты только в черновик"))
		return
	}

	if err := h.Repository.AddReagentToMethane(uint(methaneID), input.ReagentID, input.Quantity, input.OrderNum); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		// "status":  "success",
		"message": "реагент добавлен в заявку",
	})
}

func (h *Handler) UpdateMethaneReagentAPI(ctx *gin.Context) {
	methaneID, _ := strconv.Atoi(ctx.Param("id"))
	reagentID, _ := strconv.Atoi(ctx.Param("reagent_id"))

	var input struct {
		Quantity *float64 `json:"quantity"`
		OrderNum *int     `json:"order_num"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateMethaneReagent(uint(methaneID), uint(reagentID), input.Quantity, input.OrderNum); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "обновлено"})
}

func (h *Handler) RemoveReagentFromMethaneAPI(ctx *gin.Context) {
	methaneID, _ := strconv.Atoi(ctx.Param("id"))
	reagentID, _ := strconv.Atoi(ctx.Param("reagent_id"))

	if err := h.Repository.RemoveReagentFromMethane(uint(methaneID), uint(reagentID)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "реагент удалён из заявки"})
}
