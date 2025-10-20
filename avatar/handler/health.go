package handler

import (
	"context"
	"time"

	"avatar/db"
	"github.com/cloudwego/hertz/pkg/app"
)

// HealthHandler ping 默认写库（"default"）
func HealthHandler(ctx context.Context, c *app.RequestContext) {
	conn, err := db.DefaultWriteDB.DB()
	if err != nil {
		c.String(200, "can not connect to db")
		return
	}
	ctx2, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	if err := conn.PingContext(ctx2); err != nil {
		c.String(500, "db ping failed")
		return
	}
	c.String(200, "ok")
}
