package common

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"os"
	"strings"
)

type EnvType string

const (
	Prod EnvType = "prod"
	Test EnvType = "test"
	//Dev EnvType = "dev"
)

var CurrentEnv EnvType // 全局环境变量

// InitEnv 从环境变量 ENV 读取并初始化全局 Env
func InitEnv() {
	val := strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
	switch val {
	case "prod":
		CurrentEnv = Prod
	case "test":
		CurrentEnv = Test
	default:
		panic(fmt.Sprintf("未知环境变量 ENV=%s，应为 prod/test", val))
	}
	hlog.Infof("current env is: %s", val)
}
