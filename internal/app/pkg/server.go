package app

import (
	"log"
	"net/http"
	"strconv"

	docs "lab4/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"time"
	"github.com/gin-contrib/cors"
	
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
        log.Printf("REQUEST: %s %s", c.Request.Method, c.Request.URL.Path)
        c.Next()
        log.Printf("RESPONSE: %d", c.Writer.Status())
    }) 

	r.GET("/ping/:name", a.Ping)

	r.GET("/docs/doc.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, docs.SwaggerInfo.ReadDoc())
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/docs/doc.json"),
	))

	auth := r.Group("/auth")
	{
		auth.POST("/login", a.Login)
		auth.POST("/register", a.Register)
		auth.POST("/logout", a.WithJWTAuth(), a.Logout)
	}

	api := r.Group("/api")
	{
		api.GET("/reagents", a.GetReagentsPublic)
	}

	user := r.Group("/api")
	user.Use(a.WithJWTAuth())
	{
		user.GET("/methanes", a.GetMethanes)
		user.GET("/methanes/:id", a.GetMethaneByID)
		user.POST("/methanes/draft", a.CreateDraftMethane)
		user.GET("/methanes/draft", a.GetDraftMethane)
		user.PUT("/methanes/:id/form", a.FormMethane)
	}

	researcher := r.Group("/api")
	researcher.Use(a.WithJWTAuth())
	{
		researcher.PUT("/methanes/:id/complete", a.RequireModerator(), a.CompleteMethane)
	}

	addr := a.config.ServiceHost + ":" + strconv.Itoa(a.config.ServicePort)
	if a.config.ServiceHost == "" {
		addr = ":8080"
	}

	if err := r.Run(addr); err != nil {
		log.Println(err)
	}

	log.Println("Server down")
}

type pingResp struct {
	Status string `json:"status"`
}

func (a *Application) Ping(gCtx *gin.Context) {
	name := gCtx.Param("name")
	gCtx.String(http.StatusOK, "Hello %s", name)
}
