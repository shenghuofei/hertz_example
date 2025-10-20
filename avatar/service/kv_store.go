package service

import (
	db "avatar/db"
	models "avatar/db/models"
	"avatar/request"
	"avatar/utils"
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
)

func ListKVStore(ctx context.Context, req *request.ListKVStore) (*utils.PageResult[models.KVStore], error) {
	// 获取分页参数
	page := req.Page
	pageSize := req.PageSize

	// 获取可选查询条件
	keyName := req.KeyName

	// 构建 GORM 查询对象
	dbQuery := db.DefaultWriteDB.Model(&models.KVStore{})
	if keyName != "" {
		dbQuery = dbQuery.Where("key_name LIKE ?", "%"+keyName+"%")
	}

	// 执行分页查询
	var kvs []models.KVStore
	res, err := utils.Paginate[models.KVStore](dbQuery, page, pageSize, &kvs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpsertKVStore(ctx context.Context, req *request.UpsertKVStore) error {
	data := models.KVStore{
		KeyName: req.KeyName,
		Value:   req.Value,
	}
	if err := models.DBKVStore.Upsert(data); err != nil {
		hlog.Errorf("KVStore upserted key %s failed, %s", req.KeyName, err.Error())
		return err
	}
	return nil
}

func DeleteKVStore(ctx context.Context, keyName string) error {
	if err := models.DBKVStore.Delete(keyName); err != nil {
		hlog.Errorf("KVStore delete key %s failed, %s", keyName, err.Error())
		return err
	}
	return nil
}

func GetKVStore(ctx context.Context, keyName string) (*models.KVStore, error) {
	kv, err := models.DBKVStore.GetValueByKey(keyName)
	if err != nil {
		hlog.Errorf("KVStore get key %s failed, %s", keyName, err.Error())
		return nil, err
	}
	return kv, nil
}
