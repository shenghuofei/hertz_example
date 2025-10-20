package db

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
)

// 代码需要使用到的key，且存在db里面，这里统一获取一次，后面直接使用
var Mykey string

func GetInitKey() {
	mykey, err := DBKVStore.GetValueByKey("mykey")
	if err != nil {
		hlog.Errorf("get mykey err: %v", err)
	} else {
		Mykey = mykey.Value
	}
}
