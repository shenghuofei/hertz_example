package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"runtime/debug"

	"avatar/response"
	"github.com/cloudwego/hertz/pkg/app"
)

func RecoverResponse() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if r := recover(); r != nil {
				hlog.Errorf("[Panic Recovered] %v\nStack Trace:\n%s", r, debug.Stack())
				if resp, ok := r.(response.Response[any]); ok {
					// 捕获我们定义的 Fail()
					ctx.JSON(http.StatusOK, response.Response[any]{
						Code:    resp.Code,
						Message: resp.Message,
						Data:    "",
					})
					ctx.Abort()
				} else {
					// 捕获未知 panic
					ctx.JSON(http.StatusInternalServerError, response.Response[any]{
						Code:    500,
						Message: "Internal Server Error",
						Data:    "",
					})
					ctx.Abort()
				}
			}
		}()
		ctx.Next(c)
	}
}
