package response

import (
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type Response[T any] struct {
	Code    int    `json:"code"`              // 状态码，比如 0 成功，非 0 失败
	Message string `json:"message,omitempty"` // 错误或提示信息
	Data    T      `json:"data,omitempty"`    // 成功返回的数据
}

// Fail 会触发 panic，用 recover 捕获中止 handler
func Fail(c *app.RequestContext, code int, message string) {
	panic(Response[any]{
		Code:    code,
		Message: message,
		Data:    "",
	})
}

// Success 正常返回
func Success[T any](c *app.RequestContext, data T, message string) {
	c.JSON(http.StatusOK, Response[T]{
		Code:    0,
		Message: message,
		Data:    data,
	})
}
