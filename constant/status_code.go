package constant

type StatusCode int32

const (
	ServiceError          StatusCode = 10000
	UserNotLogin                     = 10001
	UserIPChanged                    = 10002
	RequestParameterError            = 10003
	FailedToProcess                  = 10004
	NoAuthority                      = 10005
)

func (s StatusCode) String() string {
	switch s {
	case ServiceError:
		return "系统错误"
	case UserNotLogin:
		return "用户未登录"
	case UserIPChanged:
		return "用户IP地址已更改"
	case RequestParameterError:
		return "请求参数错误"
	case FailedToProcess:
		return "服务器处理失败"
	case NoAuthority:
		return "没有权限"
	default:
		return "未知错误"
	}
}
