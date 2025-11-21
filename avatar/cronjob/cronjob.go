package cronjob

import (
	"avatar/db"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redsync/redsync/v4"
	"github.com/robfig/cron/v3"
	"time"
)

var (
	c *cron.Cron
)

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

// Start 启动 Cron
func Start() {
	InitCron()
	AddTasks()
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

// AddTask 注册一个定时任务,可以多副本执行
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

// AddTask 注册一个定时任务,只能单副本执行,通过分布式锁确定只能有一个执行,lockTTL需要大于任务执行时间
func AddUniqTask(spec string, task func(), lockKey string, lockTTL int) (cron.EntryID, error) {
	InitCron() // 确保 cron 初始化

	id, err := c.AddFunc(spec, func() {
		// 创建锁（过期时间可以比任务执行最长时间稍长）
		rs := db.GetRedsync()
		mutex := rs.NewMutex(lockKey, redsync.WithExpiry(time.Duration(lockTTL)*time.Second))

		// 尝试获取锁
		if err := mutex.Lock(); err != nil {
			hlog.Infof("[cronjob] skip task, another instance holds lock: %s", lockKey)
			return
		}

		// 确保 panic 时也能释放锁
		defer func() {
			if ok, err := mutex.Unlock(); err != nil || !ok {
				hlog.Errorf("[cronjob] unlock failed for %s: %v, ok=%v", lockKey, err, ok)
			} else {
				hlog.Infof("[cronjob] lock released for %s", lockKey)
			}
			if r := recover(); r != nil {
				hlog.Errorf("[cronjob] panic recovered in task %s: %v", lockKey, r)
			}
		}()

		// 执行任务
		task()
	})

	if err != nil {
		hlog.Errorf("add uniq cron task for key: %s err: %v", lockKey, err)
		return 0, err
	}
	hlog.Infof("Added uniq cron task: %s (id=%d, lockKey=%s)", spec, id, lockKey)
	return id, nil
}
