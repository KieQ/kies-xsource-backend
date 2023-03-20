package handler

import (
	"github.com/gin-gonic/gin"
	"kies-xsource-backend/constant"
	"net/http"
)

type Response struct {
	StatusCode    int32       `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data"`
}

func OnSuccess(c *gin.Context, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	successResponse := Response{Data: data}
	c.JSON(http.StatusOK, successResponse)
}

func OnFail(c *gin.Context, statusCode constant.StatusCode) {
	failResponse := Response{StatusCode: int32(statusCode), StatusMessage: statusCode.String()}
	c.JSON(http.StatusOK, failResponse)
}

func OnFailWithMessage(c *gin.Context, statusCode constant.StatusCode, statusMessage string) {
	failResponse := Response{StatusCode: int32(statusCode), StatusMessage: statusMessage}
	c.JSON(http.StatusOK, failResponse)
}
