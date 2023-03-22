package dto

type UserLoginRequest struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type UserLoginResponse struct {
	NickName string `json:"nick_name"`
	Profile  string `json:"profile"`
}

type UserSignupRequest struct {
	Account          string `json:"account"`
	Password         string `json:"password"`
	NickName         string `json:"nick_name"`
	Gender           int32  `json:"gender"`
	Profile          string `json:"profile"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	SelfIntroduction string `json:"self_introduction"`
}

type UserSignupResponse struct {
	Account string `json:"account"`
}

type UserUpdateRequest struct {
	Account          string  `json:"account"`
	Password         *string `json:"password"`
	NickName         *string `json:"nick_name"`
	Gender           *int32  `json:"gender"`
	Profile          *string `json:"profile"`
	Phone            *string `json:"phone"`
	Email            *string `json:"email"`
	SelfIntroduction *string `json:"self_introduction"`
}

type User struct {
	Account          string `json:"account"`
	NickName         string `json:"nick_name"`
	Profile          string `json:"profile"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Gender           int32  `json:"gender"`
	SelfIntroduction string `json:"self_introduction"`
	CreateTime       int64  `json:"create_time"`
}

type UserDetailRequest struct {
	Account *string `json:"account"`
}

type UserDetailResponse struct {
	User
}

type UserListRequest struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

type UserListResponse struct {
	Users []*User `json:"users"`
}
