package db

import (
	"context"
	"errors"
	"kies-xsource-backend/model/table"
	"kies-xsource-backend/utils"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
)

/**************[Write Part]***************/


func AddVoyage(ctx context.Context, voyage *table.Voyage)(int32, error) {
	logs.CtxInfo(ctx, "[DB] added voyage=%v", utils.ToJSON(voyage))
	if voyage == nil {
		return 0, errors.New("voyage is nil")
	}
	err := db.Table(table.NameVoyage).Create(voyage).Error
	return voyage.ID, err
}

func UpdateVoyage(ctx context.Context, voyageID int32, updateData map[string]any)error {
	logs.CtxInfo(ctx, "[DB] update voyage_id=%v, data=%v", voyageID, utils.ToJSON(updateData))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	err := db.Table(table.NameVoyage).Where("id = ?", voyageID).Updates(updateData).Error
	return err
}


/**************[Read Part]***************/


func GetLatestVoyageByUserID(ctx context.Context, userID int32)(*table.Voyage, error) {
	logs.CtxInfo(ctx, "[DB] user_id=%v", userID)

	var voyage *table.Voyage
	err := db.Table(table.NameVoyage).Where("user_id", userID).Last(&voyage).Error
	if err != nil{
		return nil, err
	}
	return voyage, nil
}