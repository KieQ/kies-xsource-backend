package handler

import (
	"kies-xsource-backend/constant"
	"kies-xsource-backend/dto"
	"kies-xsource-backend/service"
	"kies-xsource-backend/utils"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	req, err := utils.BindJSON[dto.UserLoginRequest](c)
	if err != nil {
		logs.CtxWarn(c, "failed to bind json, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}

	logs.CtxInfo(c, "[ENTRY] request=%v", utils.ToJSON(req))

	resp, sc, err := service.UserLogin(c, &req)
	if err != nil {
		logs.CtxWarn(c, "failed to login, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	service.SetToken(c, resp.UserID, req.RememberMe, c.GetHeader(constant.RealIP))
	logs.CtxInfo(c, "[EXIT] response=%v", utils.ToJSON(resp))
	OnSuccess(c, resp)
}

func UserSignup(c *gin.Context) {
	req, err := utils.BindJSON[dto.UserSignupRequest](c)
	if err != nil {
		logs.CtxWarn(c, "failed to bind json, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}
	logs.CtxInfo(c, "[ENTRY] request=%v", utils.ToJSON(req))

	resp, sc, err := service.UserSignup(c, &req)
	if err != nil {
		logs.CtxWarn(c, "failed to sign up, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	logs.CtxInfo(c, "[EXIT] response=%v", utils.ToJSON(resp))
	OnSuccess(c, resp)
}

func UserLogout(c *gin.Context) {

	userID := c.GetInt64(constant.UserID)

	logs.CtxInfo(c, "[ENTRY] %v logout", userID)

	c.SetCookie(constant.Token, "", -1, "/", "", false, false)

	logs.CtxInfo(c, "[EXIT]")

	OnSuccess(c, nil)
}

func UserUpdate(c *gin.Context) {
	req, err := utils.BindJSON[dto.UserUpdateRequest](c)
	if err != nil {
		logs.CtxWarn(c, "failed to bind json, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}
	logs.CtxInfo(c, "[ENTRY] request=%v", utils.ToJSON(req))

	if req.UserID != int32(c.GetInt64(constant.UserID)) {
		logs.CtxWarn(c, "user_id is not the same, user_id in request = %v, user_id in token = %v", req.UserID, c.GetInt64(constant.UserID))
		OnFail(c, constant.StatusCodeNoAuthority)
		return
	}

	sc, err := service.UserUpdate(c, &req)
	if err != nil {
		logs.CtxWarn(c, "failed to sign up, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	logs.CtxInfo(c, "[EXIT]")

	OnSuccess(c, nil)
}

func UserDetail(c *gin.Context) {
	OnSuccess(c, nil)
}

func UserList(c *gin.Context) {
	OnSuccess(c, nil)
}
