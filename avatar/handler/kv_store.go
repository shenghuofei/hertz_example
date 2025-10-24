package handler

import (
	"avatar/request"
	"avatar/response"
	"avatar/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// HealthHandler ping 默认写库（"default"）
func ListKVStore(ctx context.Context, c *app.RequestContext) {
	req := &request.ListKVStore{}
	if err := c.BindQuery(req); err != nil {
		hlog.Errorf("ListKVStore bind param err: %+v", err)
		response.Fail(c, 1, err.Error())
		return
	}
	data, err := service.ListKVStore(ctx, req)
	if err != nil {
		hlog.Errorf("ListKVStore failed to handle request: %v", err)
		response.Fail(c, 2, err.Error())
		return
	}
	response.Success(c, data, "ok")
}

func UpsertKVStore(ctx context.Context, c *app.RequestContext) {
	var req request.UpsertKVStore
	err := c.BindJSON(&req)
	hlog.Info(req.KeyName, req.Value)
	if err != nil {
		hlog.Errorf("UpsertKVStore failed to bind request: %v", err)
		response.Fail(c, 1, err.Error())
		return
	}
	if err = service.UpsertKVStore(ctx, &req); err != nil {
		hlog.Errorf("UpsertKVStore failed to handle request: %v", err)
		response.Fail(c, 2, err.Error())
		return
	}
	response.Success(c, "success", "ok")
}

func DeleteKVStoreByName(ctx context.Context, c *app.RequestContext) {
	keyName := c.Query("key_name")
	if keyName == "" {
		response.Fail(c, 1, "key_name is empty")
		return
	}
	if err := service.DeleteKVStore(ctx, keyName); err != nil {
		hlog.Errorf("DeleteKVStore for key: %s failed: %v", keyName, err)
		response.Fail(c, 2, err.Error())
		return
	}
	response.Success(c, "success", "ok")
}

func GetKVStoreByName(ctx context.Context, c *app.RequestContext) {
	keyName := c.Query("key_name")
	if keyName == "" {
		response.Fail(c, 1, "key_name is empty")
		return
	}
	kv, err := service.GetKVStore(ctx, keyName)
	if err != nil {
		hlog.Errorf("GetKVStore for key: %s failed: %v", keyName, err)
		response.Fail(c, 2, err.Error())
		return
	}
	response.Success(c, kv, "ok")
}
