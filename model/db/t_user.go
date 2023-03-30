package db

import (
	"context"
	"errors"
	"kies-xsource-backend/model/table"
	"kies-xsource-backend/utils"
	"time"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"gorm.io/gorm"
)

/**************[Write Part]***************/

func AddUser(ctx context.Context, user *table.User) (int32, error) {
	logs.CtxInfo(ctx, "added user=%v", utils.ToJSON(user))
	if user == nil {
		return 0, errors.New("user is nil")
	}
	err := db.Table(table.NameUser).Create(user).Error
	return user.ID, err
}

func UpdateUser(ctx context.Context, userID int32, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update user_id=%v, data=%v", userID, utils.ToJSON(updateData))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(table.NameUser).Where("id = ?", userID).Updates(updateData).Error
	return err
}

func UpdateUsers(ctx context.Context, userIDs []int32, updateData map[string]interface{}) error {
	logs.CtxInfo(ctx, "update data=%v, user_ids=%v", utils.ToJSON(updateData), utils.ToJSON(userIDs))
	if updateData == nil {
		return errors.New("update data is nil")
	}
	updateData["update_time"] = time.Now()
	err := db.Table(table.NameUser).Where("id in ?", userIDs).Updates(updateData).Error
	return err
}

/**************[Read Part]***************/

func GetUserWithEmail(ctx context.Context, email string) (*table.User, error) {
	logs.CtxInfo(ctx, "email=%v", email)
	users, err := GetUsersWithCondition(ctx, map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 || users[0] == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return users[0], nil
}

func GetUserWithUserID(ctx context.Context, userID string) (*table.User, error) {
	logs.CtxInfo(ctx, "user_id=%v", userID)
	users, err := GetUsersWithUserIDs(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 || users[0] == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return users[0], nil
}

func GetUsersWithUserIDs(ctx context.Context, userIDs []string) ([]*table.User, error) {
	logs.CtxInfo(ctx, "user_ids=%v", utils.ToJSON(userIDs))
	result := make([]*table.User, 0, len(userIDs))
	err := db.Table(table.NameUser).Where("id in ?", userIDs).Find(&result).Error
	return result, err
}

func GetUsersWithCondition(ctx context.Context, where map[string]interface{}) ([]*table.User, error) {
	logs.CtxInfo(ctx, "where condition=%v", utils.ToJSON(where))
	result := make([]*table.User, 0, 1)
	err := db.Table(table.NameUser).Where(where).Find(&result).Error
	return result, err
}
