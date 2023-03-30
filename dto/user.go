package dto

type UserLoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type UserLoginResponse struct {
	UserID   int32  `json:"user_id"`
	NickName string `json:"nick_name"`
	Profile  string `json:"profile"`
}

type UserSignupRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	NickName         string `json:"nick_name"`
	Gender           int8   `json:"gender"`
	Profile          string `json:"profile"`
	Phone            string `json:"phone"`
	SelfIntroduction string `json:"self_introduction"`
}

type UserSignupResponse struct {
	UserID int32 `json:"user_id"`
}

type UserUpdateRequest struct {
	UserID           int32   `json:"user_id"`
	Password         *string `json:"password"`
	NickName         *string `json:"nick_name"`
	Gender           *int8   `json:"gender"`
	Profile          *string `json:"profile"`
	Phone            *string `json:"phone"`
	SelfIntroduction *string `json:"self_introduction"`
}

type User struct {
	UserID           int32   `json:"user_id"`
	NickName         string  `json:"nick_name"`
	Profile          string  `json:"profile"`
	Phone            *string `json:"phone"`
	Email            *string `json:"email"`
	Gender           int8    `json:"gender"`
	SelfIntroduction string  `json:"self_introduction"`
	CreateTime       int64   `json:"create_time"`
}

type UserDetailRequest struct {
	UserID int32 `json:"user_id"`
}

type UserDetailResponse struct {
	User
}

type UserListRequest struct {
	Page int32 `json:"page"`
	Size int32 `json:"size"`
}

type UserListResponse struct {
	Users []*User `json:"users"`
}
