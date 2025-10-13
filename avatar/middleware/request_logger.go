package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// RequestLogger 返回一个中间件
// printBody: 是否打印请求 body
func RequestLogger(printBody bool) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		method := string(c.Method())
		path := string(c.Path())
		clientIP := c.RemoteAddr().String()
		query := string(c.Request.URI().QueryString())

		var reqBody string
		if printBody {
			body := c.Request.Body()
			if len(body) > 0 {
				// 保存原始 body
				reqBody = string(body)
				// 重新设置 body，保证后续 handler 能读取
				c.Request.SetBody(body)
			}
		}

		// 执行后续 handler
		c.Next(ctx)

		status := c.Response.StatusCode()
		respLen := len(c.Response.Body())
		cost := time.Since(start)

		hlog.Infof("[%s] %s %s from %s, query=%s body=%s, status=%d, resp_len=%d, cost=%s",
			start.Format("2006-01-02 15:04:05"),
			method,
			path,
			clientIP,
			query,
			reqBody,
			status,
			respLen,
			cost,
		)
	}
}
