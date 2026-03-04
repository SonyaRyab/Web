package handler

import (
	"lab2/internal/app/repository"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/reagents", h.GetAllReagents)
	router.GET("/reagents/:id", h.GetReagentById)
	router.GET("/methane", h.GetCurrentMethane)
	// 1 POST через ORM (добавление в заявку)
	router.POST("/add-to-methane", h.AddToMethane)
	// 1 POST через SQL UPDATE (удаление заявки)
	router.POST("/delete-methane", h.DeleteMethane)
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *Handler) UpdateMethaneItem(ctx *gin.Context) {
	appID := ctx.Param("id")

	itemIDStr := ctx.PostForm("item_id")
	quantityStr := ctx.PostForm("quantity")
	comment := ctx.PostForm("comment")

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
		return
	}

	err = h.Repository.UpdateMethaneItem(appID, itemID, quantity, comment)
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
}

// Удаление услуги из заявки
func (h *Handler) RemoveFromMethane(ctx *gin.Context) {
	appID := ctx.Param("app_id")
	itemIDStr := ctx.Param("item_id")

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/application/"+appID)
		return
	}

	err = h.Repository.RemoveFromMethane(appID, itemID)
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/application/"+appID)
}

// Очистка заявки (удаление всех услуг)
func (h *Handler) ClearApplication(ctx *gin.Context) {
	appID := ctx.Param("id")

	err := h.Repository.ClearMethane(appID)
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
}
