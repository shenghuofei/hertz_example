package db

import (
	"avatar/db"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Logs struct {
	// 如果结构体里的字段名与db里的不一样，使用column知道字段名
	ID      uint      `gorm:"primaryKey;column:id" json:"id"`
	Cluster string    `json:"cluster"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Rc      int       `json:"rc"`
	Content string    `json:"content"`
	Mods    string    `json:"mods"`
	Hosts   string    `json:"hosts"`
	// 结构体里有此字段，但是db里没有，使用-忽略此字段
	DbNoThisField string `gorm:"-" json:"db_no_this_field"`
}

var DBK8sIACLogs = &Logs{}

func (log *Logs) GetLogByhost(host string) ([]Logs, error) {
	var results []Logs
	logdb, err := db.Mgr.GetWriteDB("k8s-iac")
	if err != nil {
		return results, err
	}

	err = logdb.Where("hosts = ?", host).Find(&results).Error
	// err := db.DefaultWriteDB.Where("hosts = ?", host).Find(&results).Error
	if err != nil {
		return results, err
	}
	return results, nil
}
