package main

import (
	"fmt"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-xsource-backend/handler"
	"os"
)

func main() {
	StartServer()
}

func StartServer() {
	var port = os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

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
	user.GET("/detail", handler.UserDetail)
	user.GET("/list", handler.UserList)

	byeBye := g.Group("/bye_bye")
	byeBye.Use(handler.MiddlewareAuthority())
	byeBye.POST("/create", handler.ByeByeCreate)
	byeBye.GET("/check_result", handler.ByeByeCheckResult)
	byeBye.GET("/next_question", handler.ByeByeNextQuestion)
}
