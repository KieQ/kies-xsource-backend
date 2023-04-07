package main

import (
	"fmt"
	"kies-xsource-backend/handler"
	"kies-xsource-backend/model/db"
	"os"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
)

func main() {
	db.MustInit()
	StartServer()
}

func StartServer() {
	var port = os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	logs.Default().SetPathLength(2)
	logs.Info("Server start with port %v", port)
	server := gin.New()
	Register(server)
	if err := server.Run(fmt.Sprintf(":%v", port)); err != nil {
		panic(err)
	}
}

func Register(g *gin.Engine) {
	g.Use(gin.Logger(), gin.Recovery(), handler.MiddlewareMetaInfo())

	g.GET("/ping", handler.Ping)

	user := g.Group("/user")
	user.POST("/login", handler.UserLogin)
	user.POST("/signup", handler.UserSignup)
	user.POST("/logout", handler.MiddlewareAuthority(), handler.UserLogout)
	user.POST("/update", handler.MiddlewareAuthority(), handler.UserUpdate)
	user.GET("/detail", handler.MiddlewareAuthority(), handler.UserDetail)
	user.GET("/list", handler.MiddlewareAuthority(), handler.UserList)

	afterSale := g.Group("/after_sale")
	afterSale.Use(handler.MiddlewareAuthority())
	afterSale.GET("/voyage/check_progress", handler.AfterSaleVoyageCheckProgress)
	afterSale.POST("/voyage/start_or_continue_trip", handler.AfterSaleVoyageStartOrContinueTrip)
	afterSale.POST("/voyage/start_over", handler.AfterSaleVoyageStartOver)
	afterSale.POST("/voyage/check_result", handler.AfterSaleVoyageCheckResult)
	afterSale.POST("/voyage/next_step", handler.AfterSaleVoyageNextStep)
	afterSale.GET("/voyage/final_reward", handler.AfterSaleVoyageFinalReward)
}
