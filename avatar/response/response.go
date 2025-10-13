package response

import "github.com/cloudwego/hertz/pkg/app"

type Response[T any] struct {
	Code    int    `json:"code"`              // 状态码，比如 0 成功，非 0 失败
	Message string `json:"message,omitempty"` // 错误或提示信息
	Data    T      `json:"data,omitempty"`    // 成功返回的数据
}

func Success[T any](c *app.RequestContext, data T, message string) {
	c.JSON(200, Response[T]{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func Fail(c *app.RequestContext, code int, message string) {
	c.JSON(200, Response[any]{
		Code:    code,
		Message: message,
	})
}
