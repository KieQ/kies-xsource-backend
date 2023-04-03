package db

import (
	"context"
	"errors"
	"kies-xsource-backend/model/table"
	"kies-xsource-backend/utils"
	"time"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**************[Write Part]***************/

func AddUser(ctx context.Context, user *table.User) (int32, error) {
	logs.CtxInfo(ctx, "[DB] added user=%v", utils.ToJSON(user))
	if user == nil {
		return 0, errors.New("user is nil")
	}
	err := db.Table(table.NameUser).Create(user).Error
	return user.ID, err
}

func UpdateUser(ctx context.Context, userID int32, updateData map[string]any) error {
	logs.CtxInfo(ctx, "[DB] update user_id=%v, data=%v", userID, utils.ToJSON(updateData))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(table.NameUser).Where("id = ?", userID).Updates(updateData).Error
	return err
}

func UpdateUsers(ctx context.Context, userIDs []int32, updateData map[string]any) error {
	logs.CtxInfo(ctx, "[DB] update data=%v, user_ids=%v", utils.ToJSON(updateData), utils.ToJSON(userIDs))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(table.NameUser).Where("id in ?", userIDs).Updates(updateData).Error
	return err
}

/**************[Read Part]***************/

func GetUserWithEmail(ctx context.Context, email string) (*table.User, error) {
	logs.CtxInfo(ctx, "[DB] email=%v", email)
	users, err := GetUsersWithCondition(ctx, map[string]any{"email": email})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 || users[0] == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return users[0], nil
}

func GetUserWithUserID(ctx context.Context, userID int32) (*table.User, error) {
	logs.CtxInfo(ctx, "[DB] user_id=%v", userID)
	users, err := GetUsersWithUserIDs(ctx, []int32{userID})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 || users[0] == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return users[0], nil
}

func GetUsersWithUserIDs(ctx context.Context, userIDs []int32) ([]*table.User, error) {
	logs.CtxInfo(ctx, "[DB] user_ids=%v", utils.ToJSON(userIDs))
	result := make([]*table.User, 0, len(userIDs))
	err := db.Table(table.NameUser).Where("id in ?", userIDs).Find(&result).Error
	return result, err
}

func GetUsersWithCondition(ctx context.Context, condition map[string]any) ([]*table.User, error) {
	logs.CtxInfo(ctx, "[DB] condition=%v", utils.ToJSON(condition))
	result := make([]*table.User, 0, 1)
	err := db.Table(table.NameUser).Where(condition).Find(&result).Error
	return result, err
}

func BatchGetUsersWithOrderAndOffset(ctx context.Context, condition map[string]any, byKey string, desc bool, offset, limit int32) ([]*table.User, error) {
	logs.CtxInfo(ctx, "[DB] condition=%v, byKey=%v, desc=%v, offset=%v, limit=%v", utils.ToJSON(condition), byKey, desc, offset, limit)
	result := make([]*table.User, 0, limit)
	err := db.Table(table.NameUser).Where(condition).Order(clause.OrderByColumn{Column: clause.Column{Name: byKey}, Desc: desc}).Offset(int(offset)).Limit(int(limit)).Find(&result).Error
	return result, err
}

func GetUserCountWithCondition(ctx context.Context, condition map[string]any) (int32, error) {
	logs.CtxInfo(ctx, "[DB] condition=%v", utils.ToJSON(condition))
	var count int64
	err := db.Table(table.NameUser).Where(condition).Count(&count).Error
	return int32(count), err
}
