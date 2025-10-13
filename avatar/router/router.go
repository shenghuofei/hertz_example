package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"avatar/handler"
)

// Register 将路由注册到 hertz server，并把 db manager 注入到 handler（或用闭包）
func Register(h *server.Hertz) {
	// 健康检查路由，handler 使用 mgr
	// === health 组 ===
	health := h.Group("/healthz")
	{
		health.GET("/", func(ctx context.Context, c *app.RequestContext) {
			handler.HealthHandler(ctx, c)
		})
	}

	// === api v1 组 ===
	v1 := h.Group("/api/v1")
	{
		// 示例页面路由
		v1.GET("/", func(ctx context.Context, c *app.RequestContext) {
			handler.ExampleHandler(ctx, c)
		})
		v1.GET("/k8s/iac/log", func(ctx context.Context, c *app.RequestContext) {
			handler.GetK8sIacLogByHost(ctx, c)
		})
	}
}
