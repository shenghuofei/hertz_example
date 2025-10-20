package router

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"avatar/handler"
)

// Register 将路由注册到 hertz server
func Register(h *server.Hertz) {
	// 健康检查路由，handler 使用 mgr
	// === health 组 ===
	health := h.Group("/healthz")
	{
		health.GET("", func(ctx context.Context, c *app.RequestContext) { handler.HealthHandler(ctx, c) })
	}

	// === api v1 组 ===
	v1 := h.Group("/api/v1")
	{
		// 示例页面路由
		v1.GET("", func(ctx context.Context, c *app.RequestContext) { handler.ExampleHandler(ctx, c) })
		v1.GET("/k8s/iac/log", func(ctx context.Context, c *app.RequestContext) { handler.GetK8sIacLogByHost(ctx, c) })

		// for ky_store
		v1.GET("/kv/store/list", func(ctx context.Context, c *app.RequestContext) { handler.ListKVStore(ctx, c) })
		v1.POST("/kv/store/create", func(ctx context.Context, c *app.RequestContext) { handler.UpsertKVStore(ctx, c) })
		v1.PUT("/kv/store/update", func(ctx context.Context, c *app.RequestContext) { handler.UpsertKVStore(ctx, c) })
		v1.DELETE("/kv/store/delete", func(ctx context.Context, c *app.RequestContext) { handler.DeleteKVStoreByName(ctx, c) })
		v1.GET("/kv/store", func(ctx context.Context, c *app.RequestContext) { handler.GetKVStoreByName(ctx, c) })
	}

}
