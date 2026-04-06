package handler

import (
	"lab3/internal/app/repository"

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

func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Reagents (услуги)
	router.GET("/api/reagents", h.GetReagentsAPI)
	router.GET("/api/reagents/:id", h.GetReagentByIDAPI)
	router.POST("/api/reagents", h.AddReagentAPI)
	router.PUT("/api/reagents/:id", h.UpdateReagentAPI)
	router.DELETE("/api/reagents/:id", h.DeleteReagentAPI)

	// MethaneReagent (м-м связь)
	router.POST("/api/methanes/:id/reagents", h.AddReagentToMethaneAPI)
	router.PUT("/api/methanes/:id/reagents/:reagent_id", h.UpdateMethaneReagentAPI)
	router.DELETE("/api/methanes/:id/reagents/:reagent_id", h.RemoveReagentFromMethaneAPI)

	// Methane (заявки)
	router.GET("/api/methanes/cart", h.GetCartAPI)
	router.GET("/api/methanes", h.GetMethanesAPI)
	router.GET("/api/methanes/:id", h.GetMethaneByIdAPI)
	//router.POST("/api/methanes/:id", h.AddMethaneAPI)
	router.POST("/api/methanes", h.CreateDraftMethaneAPI)
	router.PUT("/api/methanes/:id", h.UpdateMethaneAPI)
	router.PUT("/api/methanes/:id/form", h.FormMethaneAPI)
	router.PUT("/api/methanes/:id/complete", h.CompleteMethaneAPI)
	router.DELETE("/api/methanes/:id", h.DeleteMethaneAPI)

	// Users
	router.POST("/api/users/register", h.RegisterUserAPI)
	router.POST("/api/users/login", h.LoginAPI)   // заглушка
	router.POST("/api/users/logout", h.LogoutAPI) // заглушка
	router.GET("/api/users/:id/methanes", h.GetUserMethanesAPI)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
