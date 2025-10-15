package service

import (
	db "avatar/db/models"
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func GetK8sIacLogByHost(ctx context.Context, host string) ([]db.Logs, error) {
	res, err := db.DBK8sIACLogs.GetLogByhost(host)
	return res, err
}

func CronTaskFunc() {
	hlog.Info("任务2 执行", time.Now())
}
