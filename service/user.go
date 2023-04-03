package service

import (
	"context"
	"errors"
	"kies-xsource-backend/constant"
	"kies-xsource-backend/dto"
	"kies-xsource-backend/model/db"
	"kies-xsource-backend/model/table"
	"kies-xsource-backend/utils"
	"net/mail"
	"time"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func UserLogin(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, constant.StatusCode, error) {
	if req.Email == "" || req.Password == "" {
		return nil, constant.StatusCodeRequestParameterError, errors.New("用户名或密码为空")
	}

	user, err := db.GetUserWithEmail(ctx, req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, constant.StatusCodeRequestParameterError, errors.New("用户不存在")
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to read from database, err=%v", err)
		return nil, constant.StatusCodeServiceError, errors.New("获取用户信息失败")
	}

	if user.Password != req.Password {
		logs.CtxWarn(ctx, "password is wrong")
		return nil, constant.StatusCodeRequestParameterError, errors.New("账号或者密码错误")
	}

	resp := &dto.UserLoginResponse{
		UserID:   user.ID,
		NickName: user.NickName,
		Profile:  user.Profile,
	}
	return resp, constant.StatusCodeSuccess, nil
}

func SetToken(c *gin.Context, userID int32, rememberMe bool, ip string) {
	var maxAge = 0
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	if rememberMe {
		maxAge = int(constant.RememberMeDuration.Seconds())
		claims["exp"] = time.Now().Add(constant.RememberMeDuration).Unix()
	}
	claims[constant.UserID] = userID
	claims[constant.TokenIP] = ip

	s, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logs.CtxWarn(c, "failed to signed the jwt, err", err)
		return
	}
	c.SetCookie(constant.Token, s, maxAge, "/", "", false, false)
}

func UserSignup(ctx context.Context, req *dto.UserSignupRequest) (*dto.UserSignupResponse, constant.StatusCode, error) {
	if req.Email == "" || req.Password == "" {
		return nil, constant.StatusCodeRequestParameterError, errors.New("用户邮箱或密码为空")
	}
	if req.NickName == "" {
		return nil, constant.StatusCodeRequestParameterError, errors.New("用户昵称为空")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		logs.CtxWarn(ctx, "email address is wrong, err=%v", err)
		return nil, constant.StatusCodeRequestParameterError, errors.New("邮箱不符合规范")
	}

	if _, err := db.GetUserWithEmail(ctx, req.Email); err == nil {
		return nil, constant.StatusCodeRequestParameterError, errors.New("该账户已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logs.CtxWarn(ctx, "failed to read data, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("读取数据库失败")
	}

	user := &table.User{
		Password:         req.Password,
		NickName:         req.NickName,
		Profile:          req.Profile,
		Phone:            req.Phone,
		Email:            req.Email,
		Gender:           table.UserGender(req.Gender),
		SelfIntroduction: req.SelfIntroduction,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}

	userID, err := db.AddUser(ctx, user)
	if err != nil {
		logs.CtxWarn(ctx, "failed to write data, err=%v", err)
		return nil, constant.StatusCodeServiceError, errors.New("创建用户信息失败，请重试")
	}
	resp := &dto.UserSignupResponse{
		UserID: userID,
	}
	return resp, constant.StatusCodeSuccess, nil
}

func UserUpdate(ctx context.Context, req *dto.UserUpdateRequest) (constant.StatusCode, error) {
	updateData := make(map[string]any)
	utils.AddToMapIfNotNil(updateData, req.Password, "password")
	utils.AddToMapIfNotNil(updateData, req.Phone, "phone")
	utils.AddToMapIfNotNil(updateData, req.Gender, "gender")
	utils.AddToMapIfNotNil(updateData, req.SelfIntroduction, "self_introduction")
	utils.AddToMapIfNotNil(updateData, req.Profile, "profile")
	utils.AddToMapIfNotNil(updateData, req.NickName, "nick_name")

	if len(updateData) == 0 {
		logs.CtxWarn(ctx, "nothing to update")
		return constant.StatusCodeRequestParameterError, errors.New("未传入需要更新的数据")
	}

	if err := db.UpdateUser(ctx, req.UserID, updateData); err != nil {
		logs.CtxWarn(ctx, "failed to write data, err=%v", err)
		return constant.StatusCodeFailedToProcess, errors.New("更新数据库失败")
	}
	return constant.StatusCodeSuccess, nil
}

func UserDetail(ctx context.Context, reqUserID int32) (*dto.UserDetailResponse, constant.StatusCode, error) {
	user, err := db.GetUserWithUserID(ctx, reqUserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, constant.StatusCodeRequestParameterError, errors.New("未找到用户")
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to fetch data, err=%v", err)
		return nil, constant.StatusCodeFailedToProcess, errors.New("获取用户信息失败")
	}

	userID, ok := ctx.Value(constant.UserID).(int64)
	if !ok {
		logs.CtxWarn(ctx, "user_id in context is not int64")
		return nil, constant.StatusCodeServiceError, errors.New("系统错误，请重试")
	}

	var resp = new(dto.UserDetailResponse)
	if int32(userID) == reqUserID {
		resp.UserID = user.ID
		resp.NickName = user.NickName
		resp.Profile = user.Profile
		resp.Phone = &user.Phone
		resp.Email = &user.Email
		resp.Gender = int8(user.Gender)
		resp.SelfIntroduction = user.SelfIntroduction
		resp.CreateTime = user.CreateTime.Unix()
	} else {
		resp.UserID = user.ID
		resp.NickName = user.NickName
		resp.Profile = user.Profile
		resp.Gender = int8(user.Gender)
		resp.SelfIntroduction = user.SelfIntroduction
		resp.CreateTime = user.CreateTime.Unix()
	}

	return resp, constant.StatusCodeSuccess, nil
}

func UserList(ctx context.Context, page, size int32) (*dto.UserListResponse, constant.StatusCode, error) {

	var resp = new(dto.UserListResponse)

	eg := new(errgroup.Group)
	eg.Go(func() error {
		users, err := db.BatchGetUsersWithOrderAndOffset(ctx, nil, "id", true, page*size, size)
		if err != nil {
			logs.CtxWarn(ctx, "failed to fetch data, err=%v", err)
			return errors.New("获取用户信息失败")
		}
		for _, user := range users {
			resp.Users = append(resp.Users, &dto.User{
				UserID:           user.ID,
				NickName:         user.NickName,
				Profile:          user.Profile,
				Gender:           int8(user.Gender),
				SelfIntroduction: user.SelfIntroduction,
				CreateTime:       user.CreateTime.Unix(),
			})
		}
		return nil
	})
	eg.Go(func() error {
		count, err := db.GetUserCountWithCondition(ctx, nil)
		if err != nil {
			logs.CtxWarn(ctx, "failed to fetch data, err=%v", err)
			return errors.New("获取用户数量失败")
		}
		resp.Page = page
		resp.Size = size
		resp.Total = count
		return nil
	})

	err := eg.Wait()
	if err != nil {
		return nil, constant.StatusCodeFailedToProcess, err
	}

	return resp, constant.StatusCodeSuccess, nil
}
