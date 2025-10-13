package db

import (
	"avatar/db"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Logs struct {
	ID      uint
	Cluster string
	Start   time.Time
	End     time.Time
	Rc      int
	Content string
	Mods    string
	Hosts   string
}

var DBK8sIACLogs = &Logs{}

func (log *Logs) GetLogByhost(host string) ([]Logs, error) {
	var results []Logs
	logdb, err := db.Mgr.GetWriteDB("k8s-iac")
	if err != nil {
		return results, err
	}

	err = logdb.Where("hosts = ?", host).Find(&results).Error
	if err != nil {
		return results, err
	}
	return results, nil
}
