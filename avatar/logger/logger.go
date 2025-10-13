package logger

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func Init(logFile string, level string) {
	// 使用 lumberjack 支持日志滚动
	writerLog := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // MB
		MaxBackups: 5,
		MaxAge:     7, // days
		Compress:   true,
	}

	// 同时输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, writerLog)

	// 输出到文件
	//hlog.SetOutput(writerLog)
	hlog.SetOutput(multiWriter)

	// 设置日志等级
	switch level {
	case "debug":
		hlog.SetLevel(hlog.LevelDebug)
	case "info":
		hlog.SetLevel(hlog.LevelInfo)
	case "warn":
		hlog.SetLevel(hlog.LevelWarn)
	case "error":
		hlog.SetLevel(hlog.LevelError)
	default:
		hlog.SetLevel(hlog.LevelInfo)
	}

	hlog.Infof("Logger initialized, file=%s level=%s", logFile, level)
}

func GetLogger() hlog.FullLogger {
	return hlog.DefaultLogger()
}
