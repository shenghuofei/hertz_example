package handler

import (
	"avatar/response"
	"avatar/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

// HealthHandler ping 默认写库（"default"）
func ExampleHandler(ctx context.Context, c *app.RequestContext) {
	data := map[string]string{"data": "Hertz MySQL template - OK"}
	response.Success(c, data, "ok")
}

func GetK8sIacLogByHost(ctx context.Context, c *app.RequestContext) {
	host := c.Query("host")
	res, err := service.GetK8sIacLogByHost(ctx, host)
	if err != nil {
		response.Fail(c, 1, err.Error())
	}
	response.Success(c, res, "ok")
}
