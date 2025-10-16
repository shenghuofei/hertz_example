package cronjob

import (
	"avatar/service"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

// AddTasks 添加一些示例任务,后续添加任务，只需要在这个方法里加AddTask就行了
func AddTasks() {
	// 每30秒执行一次
	AddTask("*/30 * * * * *", func() {
		hlog.Infof("任务1 每30秒执行一次 %v", time.Now())
	})

	// 每分钟执行一次且只能单实例执行，service包中的func可以作为cron来执行
	AddUniqTask("0 * * * * *", service.CronTaskFunc, "avatar:test:uinq:task", 70)
}
