package db

import (
	"avatar/db"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

/*
 * 代码需要使用到的key，且存在db里面，这里统一获取一次，后面直接使用
 * 比如mesh token等第三方接口的token
 * token不能放在配置文件中，因此统一放在db中，这里统一查询后缓存在内存中
 * 需要提前获取的在tokenCache中加上对应的key在db中的key_name,使用时通过GetToken获取
 */
var (
	tokenCache = map[string]string{
		"mykey": "",
	}
	tokenLock sync.RWMutex
)

// GetToken 获取token
func GetToken(key string) (string, bool) {
	tokenLock.RLock()
	defer tokenLock.RUnlock()
	v, ok := tokenCache[key]
	return v, ok
}

// GetInitKey 批量加载
func GetInitKey() {
	keys := []string{}
	for k := range tokenCache {
		keys = append(keys, k)
	}
	res := []KVStore{}
	err := db.DefaultWriteDB.Where("key_name IN ?", keys).Select("key_name, value").Find(&res).Error
	if err != nil {
		hlog.Errorf("get init key err: %v", err)
		panic(err)
	}
	tokenLock.Lock()
	defer tokenLock.Unlock()
	for _, kv := range res {
		tokenCache[kv.KeyName] = kv.Value
	}
}
