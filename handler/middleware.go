package handler

import (
	"kies-xsource-backend/constant"
	"kies-xsource-backend/service"
	"kies-xsource-backend/utils"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
)

func MiddlewareMetaInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.RequestID, c.GetHeader(constant.RequestID))
		c.Set(constant.RealIP, c.GetHeader(constant.RealIP))
		c.Header(constant.RequestID, c.GetHeader(constant.RequestID))

		//TODO delete debug code
		// c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		// c.Header("Access-Control-Allow-Credentials", "true")
	}
}

func MiddlewareAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get Token from cookie
		tokenStr, err := c.Cookie(constant.Token)
		if err != nil {
			logs.CtxWarn(c, "failed to get token, err=%v", err)
			OnFail(c, constant.StatusCodeUserNotLogin)
			c.Abort()
			return
		}

		//validate the JWT
		claims, err := service.ValidateToken(tokenStr)
		if err != nil {
			logs.CtxWarn(c, "failed to validate token, err=%v", err)
			OnFail(c, constant.StatusCodeUserNotLogin)
			c.Abort()
			return
		}

		//get user_id from JWT, if success, set with key user_id
		if val, err := utils.GetFromAnyMap[float64](claims, constant.UserID); err != nil {
			logs.CtxWarn(c, "JWT does not contain %v, err=%v", constant.UserID, err)
			OnFail(c, constant.StatusCodeUserNotLogin)
			c.Abort()
			return
		} else {
			c.Set(constant.UserID, int64(val))
		}

		//get the request ip and check the IP
		if val, err := utils.GetFromAnyMap[string](claims, constant.TokenIP); err != nil {
			logs.CtxWarn(c, "JWT does not contain %v, err=%v", constant.TokenIP, err)
			OnFail(c, constant.StatusCodeUserNotLogin)
			c.Abort()
			return
		} else if val != c.GetHeader(constant.RealIP) {
			logs.CtxWarn(c, "user ip has changed from %v to %v", val, c.GetHeader(constant.RealIP))
			c.SetCookie(constant.Token, "", -1, "/", "", false, false)
			OnFail(c, constant.StatusCodeUserIPChanged)
			c.Abort()
			return
		}

	}

}
