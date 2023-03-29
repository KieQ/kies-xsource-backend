package utils

import "github.com/gin-gonic/gin"

func BindJSON[T any](c *gin.Context) (T, error) {
	var req T
	err := c.ShouldBindJSON(&req)
	return req, err
}
