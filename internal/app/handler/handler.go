package handler

import (
	"lab1/internal/app/repository"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) GetReagents(ctx *gin.Context) {
	// Получаем параметр поиска из URL: /hello?query=
	var reagents []repository.Reagent
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" { // Если параметр пустой — показываем все услуги
		reagents, err = h.Repository.GetReagents()
	} else { // Иначе фильтруем по названию (поиск подстроки без учета регистра)
		reagents, err = h.Repository.GetReagentsByTitle(searchQuery)
	}
	if err != nil {
		logrus.Error(err)
	}

	apps, _ := h.Repository.GetExperiments()

	for i := range apps {
		apps[i].TotalAmount = len(apps[i].Items)
	}
	// Сохраняем query для отображения в поле поиска после запроса
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":         time.Now().Format("15:04:05"),
		"reagents":     reagents,
		"query":        searchQuery, //сохраняет поле поиска после запроса
		"applications": apps,        // Передаем заявки в шаблон
		"minioBaseUrl": "http://localhost:9000/lab1-images/",
	})
}

func (h *Handler) GetReagent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	reagent, err := h.Repository.GetReagent(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "reagent.html", gin.H{
		"reagent": reagent,
	})
}

// Страница просмотра заявки по ID
func (h *Handler) GetExperiment(ctx *gin.Context) {
	appID := ctx.Param("id") // Получаем ID заявки из URL

	application, err := h.Repository.GetExperiment(appID)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/hello")
		return
	}

	ctx.HTML(http.StatusOK, "experiment.html", gin.H{
		"time":        time.Now().Format("15:04:05"),
		"application": application,
	})
}

// Добавление услуги в заявку
func (h *Handler) AddToExperiment(ctx *gin.Context) {
	orderIDStr := ctx.PostForm("order_id")
	appID := ctx.PostForm("app_id") // Получаем ID заявки из формы

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/hello")
		return
	}

	if appID == "" {
		appID, err = h.Repository.CreateExperiment("Новый заказчик")
		if err != nil {
			logrus.Error(err)
			ctx.Redirect(http.StatusSeeOther, "/hello")
			return
		}
	}

	err = h.Repository.AddToExperiment(appID, orderID)
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
}

func (h *Handler) UpdateExperimentItem(ctx *gin.Context) {
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

	err = h.Repository.UpdateExperimentItem(appID, itemID, quantity, comment)
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
}

// Удаление услуги из заявки
func (h *Handler) RemoveFromExperiment(ctx *gin.Context) {
	appID := ctx.Param("app_id")
	itemIDStr := ctx.Param("item_id")

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusSeeOther, "/application/"+appID)
		return
	}

	err = h.Repository.RemoveFromExperiment(appID, itemID)
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/application/"+appID)
}

// Очистка заявки (удаление всех услуг)
func (h *Handler) ClearApplication(ctx *gin.Context) {
	appID := ctx.Param("id")

	err := h.Repository.ClearExperiment(appID)
	if err != nil {
		logrus.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
}
