package constant

type StatusCode int32

const (
	StatusCodeSuccess               StatusCode = 0
	StatusCodeServiceError                     = 10000
	StatusCodeUserNotLogin                     = 10001
	StatusCodeUserIPChanged                    = 10002
	StatusCodeRequestParameterError            = 10003
	StatusCodeFailedToProcess                  = 10004
	StatusCodeNoAuthority                      = 10005
)

func (s StatusCode) String() string {
	switch s {
	case StatusCodeSuccess:
		return ""
	case StatusCodeServiceError:
		return "系统错误"
	case StatusCodeUserNotLogin:
		return "用户未登录"
	case StatusCodeUserIPChanged:
		return "用户IP地址已更改"
	case StatusCodeRequestParameterError:
		return "请求参数错误"
	case StatusCodeFailedToProcess:
		return "服务器处理失败"
	case StatusCodeNoAuthority:
		return "没有权限"
	default:
		return "未知错误"
	}
}
