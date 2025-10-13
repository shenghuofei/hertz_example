package service

import (
	db "avatar/db/models"
	"context"
	_ "github.com/go-sql-driver/mysql"
)

func GetK8sIacLogByHost(ctx context.Context, host string) ([]db.Logs, error) {
	res, err := db.DBK8sIACLogs.GetLogByhost(host)
	return res, err
}
