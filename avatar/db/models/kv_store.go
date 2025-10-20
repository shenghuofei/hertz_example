package db

import (
	"avatar/db"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm/clause"
	"time"
)

type KVStore struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	KeyName   string    `gorm:"uniqueIndex;size:255;not null" json:"key_name"`
	Value     string    `gorm:"type:text;not null" json:"value"`
	UpdatedBy string    `gorm:"size:255;not null" json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var DBKVStore = &KVStore{}

func (kv *KVStore) GetValueByKey(keyName string) (*KVStore, error) {
	res := &KVStore{}
	err := db.DefaultWriteDB.Where("key_name = ?", keyName).First(res).Error
	if err != nil {
		return res, err
	}
	return res, nil
}

func (kv *KVStore) Delete(keyName string) error {
	if err := db.DefaultWriteDB.Where("key_name = ?", keyName).Delete(&KVStore{}).Error; err != nil {
		return err
	}
	return nil
}

func (kv *KVStore) Upsert(store KVStore) error {
	res := db.DefaultWriteDB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "key_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "updated_by"}),
		},
	).Create(&store)
	if res.Error != nil {
		hlog.Infof("KVStore upserted key %s failed, error: %v", store.KeyName, res.Error)
		return res.Error
	}
	hlog.Infof("KVStore upserted key %s success, effected rows: %d", store.KeyName, res.RowsAffected)
	return nil
}
