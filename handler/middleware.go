package handler

import (
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"kies-xsource-backend/constant"
	"kies-xsource-backend/service"
	"kies-xsource-backend/utils"
)

func MiddlewareMetaInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.RequestID, c.GetHeader(constant.RequestID))
		c.Set(constant.RealIP, c.GetHeader(constant.RealIP))
		s, _ := c.Cookie(constant.Language)
		c.Set(constant.Language, s)
		c.Header(constant.RequestID, c.GetHeader(constant.RequestID))

		////TODO delete debug code
		//c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		//c.Header("Access-Control-Allow-Credentials", "true")
	}
}

func MiddlewareAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get Token from cookie
		tokenStr, err := c.Cookie(constant.Token)
		if err != nil {
			logs.CtxWarn(c, "failed to get token, err=%v", err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		}

		//validate the JWT
		claims, err := service.ValidateToken(tokenStr)
		if err != nil {
			logs.CtxWarn(c, "failed to validate token, err=%v", err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		}

		//get account from JWT, if success, set with key account
		if val, err := utils.GetFromAnyMap[string](claims, constant.Account); err != nil {
			logs.CtxWarn(c, "JWT does not contain %v, err=%v", constant.Account, err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		} else {
			c.Set(constant.Account, val)
		}

		//get the request ip and check the IP
		if val, err := utils.GetFromAnyMap[string](claims, constant.TokenIP); err != nil {
			logs.CtxWarn(c, "JWT does not contain %v, err=%v", constant.TokenIP, err)
			OnFail(c, constant.UserNotLogin)
			c.Abort()
			return
		} else if val != c.GetHeader(constant.RealIP) {
			logs.CtxWarn(c, "user ip has changed from %v to %v", val, c.GetHeader(constant.RealIP))
			c.SetCookie(constant.Token, "", -1, "/", "", false, false)
			OnFail(c, constant.UserIPChanged)
			c.Abort()
			return
		}

	}

}
