package util

var (
	OK                  = &Errno{Code: 200, Message: "OK"}
	ErrAll              = &Errno{Code: 4000, Message: "通用错误信息"}
	NeedRedirect        = &Errno{Code: 3000, Message: "需要跳转"}
	InternalServerError = &Errno{Code: 1001, Message: "Internal server error"}

	ErrUserNotFound = &Errno{Code: 2102, Message: "用户未找到"}

	ErrTokenInvalid = &Errno{Code: 401, Message: "token错误"}
)
