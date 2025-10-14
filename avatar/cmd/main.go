package main

import (
	"avatar/common"
	"avatar/config"
	"avatar/db"
	"avatar/logger"
	"avatar/middleware"
	"avatar/router"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func main() {
	// get env info,must set ENV info
	common.InitEnv()
	// load config
	cfgPath := "."
	err := config.LoadConfig(cfgPath)
	if err != nil {
		msg := fmt.Sprintf("load config error: %v", err)
		panic(msg)
	}

	// init logger
	logger.Init(config.Cfg.GetString("app.log_file"), config.Cfg.GetString("app.log_level"))
	hlog.Infof("env: %s, log_level: %s", string(common.Env), config.Cfg.GetString("app.log_level"))

	// build DB manager
	// hlog.Infof("dbs :%v", config.Cfg.Sub("database").AllSettings())
	db.InitDB()

	// start prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		metric_addr := fmt.Sprintf(":%d", config.Cfg.GetInt("app.metric_port"))
		hlog.Infof("metrics server on %s", metric_addr)
		if err = http.ListenAndServe(metric_addr, nil); err != nil && err != http.ErrServerClosed {
			hlog.Fatalf("metrics server err: %v", err)
		}
	}()

	// create hertz server
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", config.Cfg.GetInt("app.port"))),
		// 优雅退出最大时长
		server.WithExitWaitTime(5*time.Second),
	)
	h.Use(middleware.RequestLogger(config.Cfg.GetBool("app.print_request_body"))) // 参数确定是否打印 body
	h.Use(middleware.RecoverResponse())                                           // 捕获handler异常

	// get default write db,global can use DefaultWriteDB value
	_, err = db.Mgr.GetDefaultWriteDB()
	if err != nil {
		hlog.Fatalf("get default write db err: %v", err)
	}

	// register routes (router is separate)
	router.Register(h)

	// 注册退出前清理逻辑
	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		hlog.Info("closing db connections...")
		db.Mgr.CloseAll()
		hlog.Info("cleanup done.")
	})

	hlog.Infof("server started on :%d", config.Cfg.GetInt("app.port"))
	h.Spin()

	// start server and graceful shutdown by self
	//go func() {
	//	if err := h.Run(); err != nil {
	//		hlog.Fatalf("server run err: %v", err)
	//	}
	//}()
	//
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//hlog.Info("shutting down...")
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//_ = h.Shutdown(ctx)
	//mgr.CloseAll()
	//hlog.Infof("shutdown complete")
}
