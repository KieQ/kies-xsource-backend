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
	"gorm.io/gorm"
)

func UserLogin(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, constant.StatusCode, error) {
	if req.Email == "" || req.Password == "" {
		return nil, constant.StatusCodeRequestParameterError, errors.New("用户名或密码为空")
	}

	user, err := db.GetUserWithEmail(ctx, req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logs.CtxWarn(ctx, "user with email %v doesn't exist, err=%v", req.Email, err)
		return nil, constant.StatusCodeRequestParameterError, errors.New("用户不存在")
	} else if err != nil {
		logs.CtxWarn(ctx, "failed to read from database, err=%v", err)
		return nil, constant.StatusCodeServiceError, errors.New("获取用户信息失败")
	}

	if user.Password != req.Password {
		logs.CtxWarn(ctx, "password is wrong")
		return nil, constant.StatusCodeRequestParameterError, errors.New("账号或者密码错误，请检查后重试")
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
		logs.CtxWarn(ctx, "failed to write to database, err=%v", err)
		return nil, constant.StatusCodeServiceError, errors.New("创建用户信息失败，请重试")
	}
	resp := &dto.UserSignupResponse{
		UserID: userID,
	}
	return resp, constant.StatusCodeSuccess, nil
}

func UserUpdate(ctx context.Context, req *dto.UserUpdateRequest) (constant.StatusCode, error) {	
	updateData := make(map[string]interface{})
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
		return constant.StatusCodeFailedToProcess, errors.New("更新数据库失败")
	}
	return constant.StatusCodeSuccess, nil
}
