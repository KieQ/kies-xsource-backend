package handler

import (
	"kies-xsource-backend/constant"
	"kies-xsource-backend/dto"
	"kies-xsource-backend/service"
	"kies-xsource-backend/utils"

	"github.com/Kidsunbo/kie_toolbox_go/cast"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
)

func AfterSaleVoyageCheckProgress(c *gin.Context) {
	userID, err := cast.To[int32](c.Query("user_id"))
	if err != nil {
		logs.CtxWarn(c, "failed to get query parameter, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}
	logs.CtxInfo(c, "[ENTRY] userID=%v", userID)

	resp, sc, err := service.AfterSaleVoyageCheckProgress(c, userID)
	if err != nil {
		logs.CtxWarn(c, "failed to check progress, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	logs.CtxInfo(c, "[EXIT] response=%v", utils.ToJSON(resp))
	OnSuccess(c, resp)
}


func AfterSaleVoyageStartOrContinueTrip(c *gin.Context) {
	req, err := utils.BindJSON[dto.AfterSaleVoyageStartOrContinueTripRequest](c)
	if err != nil {
		logs.CtxWarn(c, "failed to bind json, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}
	logs.CtxInfo(c, "[ENTRY] request=%v", utils.ToJSON(req))

	userID := int32(c.GetInt64(constant.UserID))
	logs.CtxInfo(c, "%v start or continue the trip", userID)


	resp, sc, err := service.AfterSaleVoyageStartOrContinueTrip(c, userID, req.Level)
	if err != nil {
		logs.CtxWarn(c, "failed to start or continue trip, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	logs.CtxInfo(c, "[EXIT] response=%v", utils.ToJSON(resp))
	OnSuccess(c, resp)
}

func AfterSaleVoyageStartOver(c *gin.Context) {
	req, err := utils.BindJSON[dto.AfterSaleVoyageStartOverRequest](c)
	if err != nil {
		logs.CtxWarn(c, "failed to bind json, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}
	logs.CtxInfo(c, "[ENTRY] request=%v", utils.ToJSON(req))

	userID := int32(c.GetInt64(constant.UserID))
	logs.CtxInfo(c, "%v start or continue the trip", userID)


	resp, sc, err := service.AfterSaleVoyageStartOver(c, userID, req.Level)
	if err != nil {
		logs.CtxWarn(c, "failed to start over trip, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	logs.CtxInfo(c, "[EXIT] response=%v", utils.ToJSON(resp))
	OnSuccess(c, resp)
}

func AfterSaleVoyageCheckResult(c *gin.Context) {
	req, err := utils.BindJSON[dto.AfterSaleVoyageCheckResultRequest](c)
	if err != nil {
		logs.CtxWarn(c, "failed to bind json, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}
	logs.CtxInfo(c, "[ENTRY] request=%v", utils.ToJSON(req))

	resp, sc, err := service.AfterSaleVoyageCheckResult(c, &req)
	if err != nil {
		logs.CtxWarn(c, "failed to check result, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	logs.CtxInfo(c, "[EXIT] response=%v", utils.ToJSON(resp))
	OnSuccess(c, resp)
}

func AfterSaleVoyageNextStep(c *gin.Context) {
	req, err := utils.BindJSON[dto.AfterSaleVoyageNextStepRequest](c)
	if err != nil {
		logs.CtxWarn(c, "failed to bind json, err=%v", err)
		OnFail(c, constant.StatusCodeRequestParameterError)
		return
	}
	logs.CtxInfo(c, "[ENTRY] request=%v", utils.ToJSON(req))

	resp, sc, err := service.AfterSaleVoyageNextStep(c, &req)
	if err != nil {
		logs.CtxWarn(c, "failed to start next step, err=%v", err)
		OnFailWithMessage(c, sc, err.Error())
		return
	}

	logs.CtxInfo(c, "[EXIT] response=%v", utils.ToJSON(resp))
	OnSuccess(c, resp)
}

func AfterSaleVoyageFinalReward(c *gin.Context) {
	OnSuccess(c, nil)
}
