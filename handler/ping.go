package handler

import "github.com/gin-gonic/gin"

func Ping(c *gin.Context) {
	OnSuccess(c, nil)
}
