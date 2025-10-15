package cronjob

import (
	"avatar/service"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/robfig/cron/v3"
	"time"
)

var c *cron.Cron

// 定义cron 使用的logger,仍然使用hlog打印日志
type CronPanicLogger struct{}

func (l CronPanicLogger) Info(msg string, keysAndValues ...interface{}) {}
func (l CronPanicLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	hlog.Errorf("Cron panic: %v, msg: %s, keys: %v", err, msg, keysAndValues)
}

// InitCron 初始化 cron
func InitCron() {
	if c != nil {
		return
	}
	// 支持秒级任务 + panic 捕获
	c = cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(CronPanicLogger{}),
		),
	)
}

// AddTask 注册一个定时任务
func AddTask(spec string, task func()) (cron.EntryID, error) {
	InitCron() // 确保 cron 初始化
	id, err := c.AddFunc(spec, task)
	if err != nil {
		hlog.Errorf("add cron task err:%v", err)
		return 0, err
	}
	hlog.Infof("Added cron task: %s (id=%d)", spec, id)
	return id, nil
}

// Start 启动 Cron
func Start() {
	InitCron()
	c.Start()
	hlog.Info("Cron started")
}

// Stop 停止 Cron
func Stop() {
	if c != nil {
		c.Stop()
		hlog.Info("Cron stopped")
	}
}

// ExampleTasks 添加一些示例任务,后续添加任务，只需要在这个方法里加AddTask就行了
func AddTasks() {
	AddTask("*/10 * * * * *", func() {
		hlog.Infof("任务1 每10秒执行一次 %v", time.Now())
	})

	AddTask("*/5 * * * * *", service.CronTaskFunc)
}
