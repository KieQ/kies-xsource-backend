package service

import (
	"context"
	"errors"
	"kies-xsource-backend/constant"
	"kies-xsource-backend/dto"
	"kies-xsource-backend/model/db"
	"kies-xsource-backend/model/table"
	"kies-xsource-backend/utils"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kidsunbo/kie_toolbox_go/cast"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"gorm.io/gorm"
)

func AfterSaleVoyageCheckProgress(ctx context.Context, userID int32) (*dto.AfterSaleVoyageCheckProgressResponse, constant.StatusCode, error) {
	resp := new(dto.AfterSaleVoyageCheckProgressResponse)
	voyage, err := db.GetLatestVoyageByUserID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Progress = dto.AfterSaleVoyageCheckProgressResultNeverStarted
		return resp, constant.StatusCodeSuccess, nil
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	resp.Level = int8(voyage.Level)
	if voyage.Level >= constant.AfterSaleVoyageLevelCount && voyage.Status == table.VoyageStatusPass{
		resp.Progress = dto.AfterSaleVoyageCheckProgressResultCleared
	} else {
		resp.Progress = dto.AfterSaleVoyageCheckProgressResultInTrip
	}

	return resp, constant.StatusCodeSuccess, nil
}

func AfterSaleVoyageStartOrContinueTrip(ctx context.Context, userID int32, level int8) (*dto.AfterSaleVoyageStartOrContinueTripResponse, constant.StatusCode, error) {
	voyage, err := db.GetLatestVoyageByUserID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if level != 0 {
			logs.CtxWarn(ctx, "can't find voyage but passed level is %v", level)
			return nil, constant.StatusCodeRequestParameterError, errors.New("未找到探索记录")
		}
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	var newLevel int8 = 1
	if voyage != nil {
		switch voyage.Status {
		case table.VoyageStatusPass:
			if voyage.Level >= constant.AfterSaleVoyageLevelCount {
				resp := &dto.AfterSaleVoyageStartOrContinueTripResponse{
					Cleared: true,
				}
				return resp, constant.StatusCodeSuccess, nil
			}
			if level != int8(voyage.Level) {
				return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
			}
			newLevel = level + 1
		case table.VoyageStatusFail:
			if level != int8(voyage.Level) {
				return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
			}
			newLevel = level
		case table.VoyageStatusInProgress:
			if level != int8(voyage.Level) {
				return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
			}
			records := []*table.VoyageRecord{}
			if voyage.Records != "" {
				records, err = cast.JSONTo[[]*table.VoyageRecord](voyage.Records)
				if err != nil {
					logs.CtxWarn(ctx, "failed to unmarshal records, err=%v", err)
					return nil, constant.StatusCodeFailedToProcess, errors.New("解析探索日志失败")
				}
			}
			records = append(records, &table.VoyageRecord{
				CreateTime: time.Now(),
				Action:     "continue",
			})
			updateData := make(map[string]any)
			updateData["records"] = utils.ToJSON(records)
			updateData["last_try_time"] = time.Now()
			err = db.UpdateVoyage(ctx, voyage.ID, updateData)
			if err != nil {
				logs.CtxWarn(ctx, "failed to update voyage, err=%v", err)
				return nil, constant.StatusCodeFailedToProcess, errors.New("更新探索日志失败")
			}
			resp := &dto.AfterSaleVoyageStartOrContinueTripResponse{
				Level: level,
				Seed:  voyage.Seed,
			}
			return resp, constant.StatusCodeSuccess, nil
		default:
		}
	}

	seed := rand.Int31()
	newVoyage := &table.Voyage{
		UserID:      userID,
		Seed:        seed,
		Level:       int32(newLevel),
		Status:      0,
		StartTime:   time.Now(),
		LastTryTime: time.Now(),
	}
	records := []*table.VoyageRecord{
		{
			CreateTime: time.Now(),
			Action:     "start",
		},
	}
	newVoyage.Records = utils.ToJSON(records)

	_, err = db.AddVoyage(ctx, newVoyage)
	if err != nil {
		logs.CtxWarn(ctx, "failed to add voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("创建探索失败")
	}
	resp := &dto.AfterSaleVoyageStartOrContinueTripResponse{
		Level: int8(newLevel),
		Seed:  seed,
	}
	return resp, constant.StatusCodeSuccess, nil
}

func AfterSaleVoyageStartOver(ctx context.Context, userID int32, level int8) (*dto.AfterSaleVoyageStartOverResponse, constant.StatusCode, error) {
	voyage, err := db.GetLatestVoyageByUserID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logs.CtxWarn(ctx, "can't find voyage but passed level is %v", level)
		return nil, constant.StatusCodeRequestParameterError, errors.New("未找到探索记录")
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	if voyage.Level != int32(level) {
		return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
	}

	if voyage.Status == table.VoyageStatusInProgress {
		records, err := cast.JSONTo[[]*table.VoyageRecord](voyage.Records)
		if err != nil {
			logs.CtxWarn(ctx, "failed to parse records, err=%v", err)
			return nil, constant.StatusCodeFailedToProcess, errors.New("解析探索日志失败")
		}
		records = append(records, &table.VoyageRecord{
			CreateTime: time.Now(),
			Action:     "start_over",
		})

		err = db.UpdateVoyage(ctx, voyage.ID, map[string]any{
			"status":  table.VoyageStatusFail,
			"records": utils.ToJSON(records),
		})
		if err != nil {
			logs.CtxWarn(ctx, "failed to update voyage, err=%v", err)
			return nil, constant.StatusCodeFailedToProcess, errors.New("更新当前探索进度失败")
		}
	}

	seed := rand.Int31()
	newVoyage := &table.Voyage{
		UserID:      userID,
		Seed:        seed,
		Level:       1,
		Status:      0,
		StartTime:   time.Now(),
		LastTryTime: time.Now(),
	}
	records := []*table.VoyageRecord{
		{
			CreateTime: time.Now(),
			Action:     "start",
		},
	}
	newVoyage.Records = utils.ToJSON(records)

	_, err = db.AddVoyage(ctx, newVoyage)
	if err != nil {
		logs.CtxWarn(ctx, "failed to add voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("创建探索失败")
	}
	resp := &dto.AfterSaleVoyageStartOverResponse{
		Level: 1,
		Seed:  seed,
	}
	return resp, constant.StatusCodeSuccess, nil
}

func AfterSaleVoyageCheckResult(ctx context.Context, req *dto.AfterSaleVoyageCheckResultRequest) (*dto.AfterSaleVoyageCheckResultResponse, constant.StatusCode, error) {
	userID, ok := ctx.Value(constant.UserID).(int64)
	if !ok {
		logs.CtxWarn(ctx, "user_id in context is not int64")
		return nil, constant.StatusCodeServiceError, errors.New("系统错误，请重试")
	}

	voyage, err := db.GetLatestVoyageByUserID(ctx, int32(userID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logs.CtxWarn(ctx, "can't find voyage but requested level is %v", req.Level)
		return nil, constant.StatusCodeRequestParameterError, errors.New("未找到探索记录")
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	if voyage.Level != int32(req.Level) {
		return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
	}

	pass := false
	failReason := ""
	switch req.Level {
	case 1:
		if req.Result == nil || req.IV == nil {
			logs.CtxWarn(ctx, "level 1 parameter is missing")
			return nil, constant.StatusCodeRequestParameterError, errors.New("参数错误")
		}
		now := time.Now().Unix()
		key := utils.GenerateKey(now)
		result, _ := utils.Decrypt(*req.Result, *req.IV, key[:])
		if strconv.FormatInt(int64(voyage.Seed), 10) == result {
			pass = true
			break
		}
		for i := int64(1); i < 10; i++ {
			keyLeft := utils.GenerateKey(now - i)
			result, _ = utils.Decrypt(*req.Result, *req.IV, keyLeft[:])
			if strconv.FormatInt(int64(voyage.Seed), 10) == result {
				pass = true
				break
			}

			keyRight := utils.GenerateKey(now + i)
			result, _ = utils.Decrypt(*req.Result, *req.IV, keyRight[:])
			if strconv.FormatInt(int64(voyage.Seed), 10) == result {
				pass = true
				break
			}
		}
		if !pass {
			failReason = "请求参数可能为伪造"
		}
	case 2:
		if req.Result == nil {
			logs.CtxWarn(ctx, "level 2 parameter is missing")
			return nil, constant.StatusCodeRequestParameterError, errors.New("参数错误")
		}
		if strings.ToLower(*req.Result) == os.Getenv("LEVEL_2_RESULT") {
			pass = true
		}
		if !pass {
			failReason = "结果不正确"
		}
	default:
		logs.CtxWarn(ctx, "unknown level value %v", req.Level)
		return nil, constant.StatusCodeRequestParameterError, errors.New("参数错误")
	}

	if pass {
		logs.CtxInfo(ctx, "user %v has passed level%v", userID, req.Level)

		records, err := cast.JSONTo[[]*table.VoyageRecord](voyage.Records)
		if err != nil {
			logs.CtxWarn(ctx, "failed to parse records, err=%v", err)
			return nil, constant.StatusCodeFailedToProcess, errors.New("解析探索日志失败")
		}
		now := time.Now()
		records = append(records, &table.VoyageRecord{
			CreateTime: now,
			Action:     "pass",
		})

		err = db.UpdateVoyage(ctx, voyage.ID, map[string]any{
			"status":    table.VoyageStatusPass,
			"records":   utils.ToJSON(records),
			"pass_time": now,
		})
		if err != nil {
			logs.CtxWarn(ctx, "failed to update voyage, err=%v", err)
			return nil, constant.StatusCodeFailedToProcess, errors.New("更新当前探索进度失败")
		}
	} else {
		logs.CtxInfo(ctx, "user %v has not passed level%v", userID, req.Level)
	}
	resp := &dto.AfterSaleVoyageCheckResultResponse{
		Pass:       pass,
		FailReason: failReason,
	}
	return resp, constant.StatusCodeSuccess, nil
}

func AfterSaleVoyageNextStep(ctx context.Context, req *dto.AfterSaleVoyageNextStepRequest) (*dto.AfterSaleVoyageNextStepResponse, constant.StatusCode, error) {
	userID, ok := ctx.Value(constant.UserID).(int64)
	if !ok {
		logs.CtxWarn(ctx, "user_id in context is not int64")
		return nil, constant.StatusCodeServiceError, errors.New("系统错误，请重试")
	}

	voyage, err := db.GetLatestVoyageByUserID(ctx, int32(userID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logs.CtxWarn(ctx, "can't find voyage but requested level is %v", req.FinishedLevel)
		return nil, constant.StatusCodeRequestParameterError, errors.New("未找到探索记录")
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	if voyage.Level != int32(req.FinishedLevel) {
		return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
	}

	if voyage.Status != table.VoyageStatusPass {
		return nil, constant.StatusCodeRequestParameterError, errors.New("当前关卡未通过，无法进入下一关卡")
	}

	records, err := cast.JSONTo[[]*table.VoyageRecord](voyage.Records)
	if err != nil {
		logs.CtxWarn(ctx, "failed to parse records, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("解析探索日志失败")
	}
	records = append(records, &table.VoyageRecord{
		CreateTime: time.Now(),
		Action:     "pass",
	})

	err = db.UpdateVoyage(ctx, voyage.ID, map[string]any{
		"status":  table.VoyageStatusPass,
		"records": utils.ToJSON(records),
		"pass_time": time.Now(),
	})
	if err != nil {
		logs.CtxWarn(ctx, "failed to update voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("更新当前探索进度失败")
	}

	if req.FinishedLevel == constant.AfterSaleVoyageLevelCount {
		logs.CtxInfo(ctx, "user %v has passed all the levels", userID)
		resp := &dto.AfterSaleVoyageNextStepResponse{
			Cleared: true,
		}
		return resp, constant.StatusCodeSuccess, nil
	}

	seed := rand.Int31()
	next := req.FinishedLevel + 1
	newVoyage := &table.Voyage{
		UserID:      int32(userID),
		Seed:        seed,
		Level:       int32(next),
		Status:      0,
		StartTime:   time.Now(),
		LastTryTime: time.Now(),
	}
	records = []*table.VoyageRecord{
		{
			CreateTime: time.Now(),
			Action:     "start",
		},
	}
	newVoyage.Records = utils.ToJSON(records)

	_, err = db.AddVoyage(ctx, newVoyage)
	if err != nil {
		logs.CtxWarn(ctx, "failed to add voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("创建探索失败")
	}
	resp := &dto.AfterSaleVoyageNextStepResponse{
		Level: next,
		Seed:  seed,
	}
	return resp, constant.StatusCodeSuccess, nil
}
