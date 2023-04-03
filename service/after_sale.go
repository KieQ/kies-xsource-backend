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
	"time"

	"github.com/Kidsunbo/kie_toolbox_go/cast"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"gorm.io/gorm"
)

func AfterSaleVoyageCheckProgress(ctx context.Context, userID int32) (*dto.AfterSaleVoyageCheckProgressResponse, constant.StatusCode, error) {
	resp := new(dto.AfterSaleVoyageCheckProgressResponse)
	resp.UserID = userID
	voyage, err := db.GetLatestVoyageByUserID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Progress = dto.AfterSaleVoyageCheckProgressResultNeverStarted
		return resp, constant.StatusCodeSuccess, nil
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	if voyage.Status == table.VoyageStatusInProgress || voyage.Status == table.VoyageStatusFail {
		resp.Progress = dto.AfterSaleVoyageCheckProgressResultInTrip
		resp.Level = int8(voyage.Level)
	} else {
		if voyage.Level >= constant.AfterSaleVoyageLevelCount {
			resp.Progress = dto.AfterSaleVoyageCheckProgressResultPass
		} else {
			resp.Progress = dto.AfterSaleVoyageCheckProgressResultInTrip
			resp.Level = int8(voyage.Level) + 1
		}
	}

	return resp, constant.StatusCodeSuccess, nil
}

func AfterSaleVoyageStartOrContinueTrip(ctx context.Context, userID int32, level int8)(*dto.AfterSaleVoyageStartOrContinueTripResponse, constant.StatusCode, error){
	voyage, err := db.GetLatestVoyageByUserID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound){
		if level != 0{
			logs.CtxWarn(ctx, "can't find voyage but passed level is %v", level)
			return nil, constant.StatusCodeRequestParameterError, errors.New("未找到探索记录")
		}
	}else if err != nil{
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	var newLevel int8 = 1
	if voyage != nil{
		switch voyage.Status{
		case table.VoyageStatusPass:
			if voyage.Level >= constant.AfterSaleVoyageLevelCount{
				resp := &dto.AfterSaleVoyageStartOrContinueTripResponse{
					Passed: true,
				}
				return resp, constant.StatusCodeSuccess, nil
			}
			if level != int8(voyage.Level) + 1{
				return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
			}
			newLevel = level
		case table.VoyageStatusFail:
			if level != int8(voyage.Level){
				return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
			}
			newLevel = level
		case table.VoyageStatusInProgress:
			if level != int8(voyage.Level){
				return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
			}
			records := []*table.VoyageRecord{}
			if voyage.Records != ""{
				records, err = cast.JSONTo[[]*table.VoyageRecord](voyage.Records)
				if err != nil{
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
			if err != nil{
				logs.CtxWarn(ctx, "failed to update voyage, err=%v", err)
				return nil, constant.StatusCodeFailedToProcess, errors.New("更新探索日志失败")
			}
			resp := &dto.AfterSaleVoyageStartOrContinueTripResponse{
				Level:  level,
			}
			return resp, constant.StatusCodeSuccess, nil
		default:
		}
	}

	newVoyage := &table.Voyage{
		UserID:      userID,
		Seed:        rand.Int31(),
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
	if err != nil{
		logs.CtxWarn(ctx, "failed to add voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("创建探索失败")
	}
	resp := &dto.AfterSaleVoyageStartOrContinueTripResponse{
		Level: int8(newLevel),
	}
	return resp, constant.StatusCodeSuccess, nil
}

func AfterSaleVoyageStartOver(ctx context.Context, userID int32, level int8)(*dto.AfterSaleVoyageStartOverResponse, constant.StatusCode, error){
	voyage, err := db.GetLatestVoyageByUserID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound){
		logs.CtxWarn(ctx, "can't find voyage but passed level is %v", level)
		return nil, constant.StatusCodeRequestParameterError, errors.New("未找到探索记录")
	}else if err != nil{
		logs.CtxWarn(ctx, "failed to fetch voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取探索进度失败")
	}

	if voyage.Level != int32(level){
		return nil, constant.StatusCodeRequestParameterError, errors.New("探索进度不匹配，请刷新后重试")
	}

	if voyage.Status == table.VoyageStatusInProgress{
		records, err := cast.JSONTo[[]*table.VoyageRecord](voyage.Records)
		if err != nil{
			logs.CtxWarn(ctx, "failed to parse records, err=%v", err)
			return nil, constant.StatusCodeFailedToProcess, errors.New("解析探索日志失败")
		}
		records = append(records, &table.VoyageRecord{
			CreateTime: time.Now(),
			Action:     "start_over",
		})

		err = db.UpdateVoyage(ctx, voyage.ID, map[string]any{
			"status":table.VoyageStatusFail,
			"records": utils.ToJSON(records),
		})
		if err != nil{
			logs.CtxWarn(ctx, "failed to update voyage, err=%v", err)
			return nil, constant.StatusCodeFailedToProcess, errors.New("更新当前探索进度失败")
		}
	}

	newVoyage := &table.Voyage{
		UserID:      userID,
		Seed:        rand.Int31(),
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
	if err != nil{
		logs.CtxWarn(ctx, "failed to add voyage, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("创建探索失败")
	}
	resp := &dto.AfterSaleVoyageStartOverResponse{
		Level: 1,
	}
	return resp, constant.StatusCodeSuccess, nil
}