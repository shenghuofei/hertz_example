package db

import (
	"avatar/db"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"strings"
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

func BuildWhere(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for col, val := range filters {
		if val == nil {
			continue
		}
		switch v := val.(type) {
		case string:
			if v == "" {
				continue
			}
			query = query.Where(fmt.Sprintf("%s = ?", col), v)
		default:
			query = query.Where(fmt.Sprintf("%s = ?", col), v)
		}
	}
	return query
}

/*
	根据map自动生成where查询条件，类似BuildWhere但更强大，使用示例

filters := map[string]interface{}{
"is_delete": 0,
"state":     []string{"online", "mounted"},
"env":       "prod",
"created_at >=": "2025-10-01",
"created_at <=": "2025-10-23",
"os_info LIKE": "%centos%",
}

query := ApplyFilters(db.Table("my_table_name"), filters)

err := query.Select(`

	SUM(cpu_core) AS cpu_core,
	SUM(memory) AS memory,
	COUNT(*) AS host_num,
	location, env, business, type,
	? AS date`, date).

Group("location, env, business, type").
Scan(&results).Error
*/
func ApplyFilters(db *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for field, val := range filters {
		// 1️⃣ slice => IN 查询
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Slice && !strings.EqualFold(field, "BETWEEN") {
			if v.Len() == 0 || (v.Len() == 1 && v.Index(0).IsZero()) {
				continue // 空切片跳过
			}
			db = db.Where(field+" IN ?", val)
			continue
		}

		// 2️⃣ BETWEEN 语法（值必须是长度为2的切片）
		if field == "BETWEEN" {
			// val = map[string][]interface{}
			if betweenMap, ok := val.(map[string][]interface{}); ok {
				for key, arr := range betweenMap {
					if len(arr) == 2 {
						db = db.Where(key+" BETWEEN ? AND ?", arr[0], arr[1])
					}
				}
			}
			continue
		}

		// 3️⃣ LIKE / >= / <= 等运算符
		if hasOp(field) {
			db = db.Where(field+" ?", val)
			continue
		}

		// 4️⃣ 默认 =
		if v.Kind() == reflect.String && v.String() == "" {
			continue // 空字符串跳过
		}
		db = db.Where(field+" = ?", val)
	}

	return db
}

func hasOp(field string) bool {
	ops := []string{"LIKE", ">=", "<=", ">", "<", "!="}
	for _, op := range ops {
		if len(field) >= len(op) && field[len(field)-len(op):] == op {
			return true
		}
	}
	return false
}
