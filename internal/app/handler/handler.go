package handler

import (
	"lab2/internal/app/repository"

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

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
//func (h *Handler) RegisterHandler(router *gin.Engine) {
//	router.GET("/hello", h.GetOrders)
//	router.GET("/order/:id", h.GetOrder)
//}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/reagents", h.GetAllReagents)
	router.GET("/reagents/:id", h.GetReagentById)
	router.GET("/experiment", h.GetCurrentExperiment)
	// 1 POST через ORM (добавление в заявку)
	router.POST("/add-to-experiment", h.AddToExperiment)
	// 1 POST через SQL UPDATE (удаление заявки)
	router.POST("/delete-experiment", h.DeleteExperiment)
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
//func (h *Handler) RegisterStatic(router *gin.Engine) {
//	router.LoadHTMLGlob("templates/*")
//	router.Static("/styles", "./styles")
//}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

// errorHandler для более удобного вывода ошибок
//func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
//	logrus.Error(err.Error())
//	ctx.JSON(errorStatusCode, gin.H{
//		"status":      "error",
//		"description": err.Error(),
//	})
//}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

//func (h *Handler) GetOrders(ctx *gin.Context) {
//	// Получаем параметр поиска из URL: /hello?query=
//	var orders []repository.Order
//	var err error
//
//	searchQuery := ctx.Query("query")
//	if searchQuery == "" { // Если параметр пустой — показываем все услуги
//		orders, err = h.Repository.GetOrders()
//	} else { // Иначе фильтруем по названию (поиск подстроки без учета регистра)
//		orders, err = h.Repository.GetOrdersByTitle(searchQuery)
//	}
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	apps, _ := h.Repository.GetExperiments()
//
//	for i := range apps {
//		apps[i].TotalAmount = len(apps[i].Items)
//	}
//	// Сохраняем query для отображения в поле поиска после запроса
//	ctx.HTML(http.StatusOK, "index.html", gin.H{
//		"time":         time.Now().Format("15:04:05"),
//		"orders":       orders,
//		"query":        searchQuery, //сохраняет поле поиска после запроса
//		"applications": apps,        // Передаем заявки в шаблон
//		"minioBaseUrl": "http://localhost:9000/lab1-images/",
//	})
//}

//func (h *Handler) GetOrder(ctx *gin.Context) {
//	idStr := ctx.Param("id")
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	order, err := h.Repository.GetOrder(id)
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	ctx.HTML(http.StatusOK, "order.html", gin.H{
//		"order": order,
//	})
//}

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
//func (h *Handler) AddToExperiment(ctx *gin.Context) {
//	orderIDStr := ctx.PostForm("order_id")
//	appID := ctx.PostForm("app_id") // Получаем ID заявки из формы
//
//	orderID, err := strconv.Atoi(orderIDStr)
//	if err != nil {
//		logrus.Error(err)
//		ctx.Redirect(http.StatusSeeOther, "/hello")
//		return
//	}
//
//	if appID == "" {
//		appID, err = h.Repository.CreateExperiment("Новый заказчик")
//		if err != nil {
//			logrus.Error(err)
//			ctx.Redirect(http.StatusSeeOther, "/hello")
//			return
//		}
//	}
//
//	err = h.Repository.AddToExperiment(appID, orderID)
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	ctx.Redirect(http.StatusSeeOther, "/application/view/"+appID)
//}

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
