package api

import (
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("C:/VSCodeProjects/WEB/lab1/templates/*")
	r.Static("/static", "C:/VSCodeProjects/WEB/lab1/resources")

	r.GET("/hello", handler.GetReagents)                 //список услуг
	r.GET("/reagent/:id", handler.GetReagent)            //детальная страница услуги
	r.GET("/experiment/view/:id", handler.GetExperiment) //страница заявки

	//r.POST("/experiment/add", handler.AddToExperiment)
	//r.POST("/experiment/:id/clear", handler.ClearExperiment)
	//r.POST("/experiment/:id/update-item", handler.UpdateExperimentItem)
	//r.GET("/experiment/:app_id/remove/:item_id", handler.RemoveFromExperiment)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
