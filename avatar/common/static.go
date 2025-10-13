package common

import (
	"fmt"
	"os"
	"strings"
)

type EnvType string

const (
	Prod EnvType = "prod"
	Test EnvType = "test"
	//Dev EnvType = "dev"
)

var Env EnvType // 全局环境变量

// InitEnv 从环境变量 ENV 读取并初始化全局 Env
func InitEnv() {
	val := strings.ToLower(strings.TrimSpace(os.Getenv("ENV")))
	switch val {
	case "prod":
		Env = Prod
	case "test":
		Env = Test
	default:
		panic(fmt.Sprintf("未知环境变量 ENV=%s，应为 prod/test", val))
	}
}
