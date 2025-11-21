//package db
//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"github.com/cloudwego/hertz/pkg/common/hlog"
//	"gorm.io/gorm"
//	"gorm.io/gorm/logger"
//	"net/url"
//	"reflect"
//	"sync"
//	"sync/atomic"
//	"time"
//
//	"avatar/config"
//	_ "github.com/go-sql-driver/mysql"
//	"gorm.io/driver/mysql"
//)
//
//// DBGroup 表示一组 DB（写 + 多个读）
//type DBGroup struct {
//	Write *gorm.DB
//	Reads []*gorm.DB
//	rr    uint32
//	mu    sync.RWMutex
//}
//
//// Manager 管理多个 DBGroup
//type Manager struct {
//	mu     sync.RWMutex
//	groups map[string]*DBGroup
//}
//
//func NewManager() *Manager {
//	return &Manager{
//		groups: make(map[string]*DBGroup),
//	}
//}
//
//func InitDBGroup(db *gorm.DB) *DBGroup {
//	// build DB manager
//	mgr := NewManager()
//	for _, name := range config.Cfg.GetStringSlice("database") {
//		if err := mgr.AddGroup(name); err != nil {
//			hlog.Fatalf("add db group %s err: %v", name, err)
//		}
//		hlog.Infof("db group %s added", name)
//	}
//}
//
//func buildDSN(c config.DBCfg, host string) string {
//	if c.DSN != "" {
//		return c.DSN
//	}
//	user := url.QueryEscape(c.User)
//	pass := url.QueryEscape(c.Password)
//
//	params := url.Values{}
//	v := reflect.ValueOf(c.Params)
//	t := reflect.TypeOf(c.Params)
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		key := field.Tag.Get("mapstructure") // 使用 mapstructure tag 作为 key
//		value := v.Field(i).String()
//		if value != "" {
//			params.Set(key, value)
//		}
//	}
//
//	//for k, v := range c.Params. {
//	//	hlog.Infof("%s=%v", k, v)
//	//	params.Set(k, v)
//	//}
//	p := ""
//	if len(params) > 0 {
//		p = "?" + params.Encode()
//	}
//	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s%s", user, pass, host, c.DBName, p)
//	//hlog.Infof("dsn: %s", dsn)
//	return dsn
//}
//
//func openDBWithPing(dsn string, pool config.DBPoolConfig, pingTimeout time.Duration) (*gorm.DB, error) {
//	// 先创建 sql.DB 测试连接
//	sqlDB, err := sql.Open("mysql", dsn)
//	if err != nil {
//		return nil, fmt.Errorf("sql.Open failed: %w", err)
//	}
//	//defer sqlDB.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
//	defer cancel()
//	if err := sqlDB.PingContext(ctx); err != nil {
//		return nil, fmt.Errorf("db ping failed: %w", err)
//	}
//
//	// 连接池设置
//	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
//	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
//	sqlDB.SetConnMaxLifetime(time.Duration(pool.ConnMaxLifetime) * time.Minute)
//
//	// GORM 打开
//	db, err := gorm.Open(mysql.New(mysql.Config{
//		Conn: sqlDB,
//	}), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info),
//	})
//	if err != nil {
//		return nil, fmt.Errorf("gorm.Open failed: %w", err)
//	}
//
//	return db, nil
//}
//
//// AddGroup 根据配置新增一个 DBGroup（会立即尝试连接并验证）
//func (m *Manager) AddGroup(name string) error {
//	group := &DBGroup{}
//
//	// build write
//	dbcfg := config.Cfg.Sub(fmt.Sprintf("database.%s", name))
//	if cfg.Hosts.Write != "" || cfg.DSN != "" {
//		dsn := buildDSN(cfg, cfg.Hosts.Write)
//		pingTimeout := time.Duration(cfg.Health.PingTimeoutMs) * time.Millisecond
//		wdb, err := openDBWithPing(dsn, cfg.Pool, pingTimeout)
//		if err != nil {
//			return fmt.Errorf("open write db %s err: %w", name, err)
//		}
//		group.Write = wdb
//	}
//
//	// build read replicas
//	for _, r := range cfg.Hosts.Reads {
//		dsn := buildDSN(cfg, r)
//		pingTimeout := time.Duration(cfg.Health.PingTimeoutMs) * time.Millisecond
//		rdb, err := openDBWithPing(dsn, cfg.Pool, pingTimeout)
//		if err != nil {
//			// cleanup
//			//for _, x := range group.Reads {
//			//	if sqlDB, e := x.DB(); e == nil {
//			//		sqlDB.Close()
//			//	}
//			//}
//			//if group.Write != nil {
//			//	if sqlDB, e := group.Write.DB(); e == nil {
//			//		sqlDB.Close()
//			//	}
//			//}
//			hlog.Errorf("open read db %s for %s err: %w", name, r, err)
//		}
//		group.Reads = append(group.Reads, rdb)
//	}
//
//	// 如果有读库，但是一个也没连上，报错
//	if len(cfg.Hosts.Reads) > 0 && len(group.Reads) == 0 {
//		return fmt.Errorf("can not open any read db for %s", name)
//	}
//
//	m.mu.Lock()
//	m.groups[name] = group
//	m.mu.Unlock()
//	return nil
//}
//
//// GetWriteDB 返回写库
//func (m *Manager) GetWriteDB(name string) (*gorm.DB, error) {
//	m.mu.RLock()
//	g, ok := m.groups[name]
//	m.mu.RUnlock()
//	if !ok || g.Write == nil {
//		return nil, fmt.Errorf("no write db for %s", name)
//	}
//	return g.Write, nil
//}
//
//// GetReadDB 返回读库（轮询），没有读库则 fallback 到写库
//func (m *Manager) GetReadDB(name string) (*gorm.DB, error) {
//	m.mu.RLock()
//	g, ok := m.groups[name]
//	m.mu.RUnlock()
//	if !ok {
//		return nil, fmt.Errorf("no db group %s", name)
//	}
//	g.mu.RLock()
//	n := len(g.Reads)
//	g.mu.RUnlock()
//	if n == 0 {
//		if g.Write != nil {
//			return g.Write, nil
//		}
//		return nil, fmt.Errorf("no read/write db available for %s", name)
//	}
//	idx := int(atomic.AddUint32(&g.rr, 1)) % n
//	g.mu.RLock()
//	db := g.Reads[idx]
//	g.mu.RUnlock()
//	return db, nil
//}
//
//func (m *Manager) CloseAll() {
//	m.mu.RLock()
//	defer m.mu.RUnlock()
//	for _, g := range m.groups {
//		if g.Write != nil {
//			if sqlDB, err := g.Write.DB(); err == nil {
//				sqlDB.Close()
//			}
//			g.Write = nil
//		}
//		for i, r := range g.Reads {
//			if sqlDB, err := r.DB(); err == nil {
//				sqlDB.Close()
//			}
//			g.Reads[i] = nil
//		}
//		g.Reads = nil
//	}
//}
//
//var DefaultWriteDB *gorm.DB
//
//func GetDefaultWriteDB(m *Manager) (*gorm.DB, error) {
//	var err error
//	DefaultWriteDB, err = m.GetWriteDB("default")
//	if err != nil {
//		return nil, err
//	}
//	return DefaultWriteDB, nil
//}

package db

import (
	"avatar/config"
	"context"
	"database/sql"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBGroup 表示一组 DB（写 + 多个读）
type DBGroup struct {
	Write *gorm.DB
	Reads []*gorm.DB
	rr    uint32
	mu    sync.RWMutex
}

// Manager 管理多个 DBGroup
type Manager struct {
	mu     sync.RWMutex
	groups map[string]*DBGroup
}

func NewManager() *Manager {
	return &Manager{
		groups: make(map[string]*DBGroup),
	}
}

var Mgr = NewManager()
var DefaultWriteDB *gorm.DB

type GormHlogWriter struct{}

func (w *GormHlogWriter) Write(p []byte) (n int, err error) {
	// 去除多余换行，避免重复打印
	msg := string(p)
	hlog.Info(msg)
	return len(p), nil
}

func InitDB() {
	databases := config.Cfg.Sub("database").AllSettings()
	for name := range databases {
		hlog.Infof("InitDB: %s", name)
		if err := Mgr.AddGroup(name); err != nil {
			hlog.Fatalf("add db group %s err: %v", name, err)
		}
		hlog.Infof("db group %s added", name)
	}
	// get default write db,global can use DefaultWriteDB value
	_, err := Mgr.GetDefaultWriteDB()
	if err != nil {
		hlog.Fatalf("get default write db err: %v", err)
	}

	// 启动健康检查，暂时先不启用
	//go Mgr.StartHealthCheck()
}

// buildDSN 根据 viper 配置构造 DSN
func buildDSN(instanceName string, host string) string {
	sub := config.Cfg.Sub(fmt.Sprintf("database.%s", instanceName))
	if sub == nil {
		return ""
	}

	// 如果配置了 DSN，直接使用
	//dsn := sub.GetString("dsn")
	//if dsn != "" {
	//	return dsn
	//}

	user := sub.GetString("user")
	//hlog.Infof("instnace %s user: %s", instanceName, user)
	pass := sub.GetString("password")
	//hlog.Infof("instnace %s pass: %s", instanceName, pass)
	db := sub.GetString("dbname")

	// params
	params := url.Values{}
	// 默认参数，优化连接稳定性
	params.Set("parseTime", "true")
	params.Set("loc", "Local")
	params.Set("charset", "utf8mb4")
	params.Set("timeout", "10s")
	params.Set("readTimeout", "30s")
	params.Set("writeTimeout", "30s")

	p := sub.Sub("params")
	// Viper会将所有的key转为小写，被转为小写以后，db初始化会报错，这里需要映射回去
	paramsMap := map[string]string{
		"parsetime":    "parseTime",
		"readtimeout":  "readTimeout",
		"writetimeout": "writeTimeout",
	}
	if p != nil {
		for _, k := range p.AllKeys() {
			if key, ok := paramsMap[k]; ok {
				params.Set(key, p.GetString(k))
			} else {
				params.Set(k, p.GetString(k))
			}
		}
	}

	paramStr := ""
	if len(params) > 0 {
		paramStr = "?" + params.Encode()
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s%s", user, pass, host, db, paramStr)
	return dsn
}

// openDBWithPing 创建 DB 并测试连接
func openDBWithPing(dsn string, dbName string) (*gorm.DB, error) {
	sub := config.Cfg.Sub(fmt.Sprintf("database.%s.pool", dbName))
	if sub == nil {
		sub = viper.New()
	}

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	pingTimeout := time.Duration(config.Cfg.GetInt(fmt.Sprintf("database.%s.health.ping_timeout_ms", dbName))) * time.Millisecond
	hlog.Infof("ping timeout: %s", pingTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	// 连接池配置
	maxOpenConns := sub.GetInt("max_open_conns")
	if maxOpenConns == 0 {
		maxOpenConns = 100 // 默认值
	}
	maxIdleConns := sub.GetInt("max_idle_conns")
	if maxIdleConns == 0 {
		maxIdleConns = 10 // 默认值
	}
	connMaxLifetime := sub.GetInt("conn_max_lifetime_min")
	if connMaxLifetime == 0 {
		connMaxLifetime = 60 // 默认60分钟
	}
	connMaxIdleTime := sub.GetInt("conn_max_idle_time_min")
	if connMaxIdleTime == 0 {
		connMaxIdleTime = 10 // 默认10分钟
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Minute)

	// 用 hlog 作为 GORM 输出
	slowMs := 1000 // 慢查询阈值 1s
	gormLogger := logger.New(
		log.New(&GormHlogWriter{}, "", 0), // 用我们上面定义的 Writer
		logger.Config{
			SlowThreshold: time.Duration(slowMs) * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}
	return db, nil
}

// AddGroupFromViper 根据 viper 配置加载数据库组
func (m *Manager) AddGroup(instanceName string) error {
	group := &DBGroup{}
	sub := config.Cfg.Sub(fmt.Sprintf("database.%s", instanceName))
	if sub == nil {
		return fmt.Errorf("database.%s config not found", instanceName)
	}

	// 如果配置了 DSN，直接使用
	dsn := sub.GetString("dsn")
	if dsn != "" {
		wdb, err := openDBWithPing(dsn, instanceName)
		if err != nil {
			return fmt.Errorf("open db %s err: %w", instanceName, err)
		}
		group.Write = wdb
		group.Reads = append(group.Reads, wdb)
		m.mu.Lock()
		m.groups[instanceName] = group
		m.mu.Unlock()
		return nil
	}

	hostsSub := config.Cfg.Sub(fmt.Sprintf("database.%s.hosts", instanceName))
	if hostsSub == nil {
		return fmt.Errorf("database.%s.hosts not found", instanceName)
	}

	writeHost := hostsSub.GetString("write")
	readHosts := hostsSub.GetStringSlice("reads")

	// 写库
	if writeHost != "" {
		dsn = buildDSN(instanceName, writeHost)
		wdb, err := openDBWithPing(dsn, instanceName)
		if err != nil {
			return fmt.Errorf("open write db %s err: %w", instanceName, err)
		}
		group.Write = wdb
	}

	// 读库
	for _, r := range readHosts {
		dsn = buildDSN(instanceName, r)
		rdb, err := openDBWithPing(dsn, instanceName)
		if err != nil {
			hlog.Errorf("open read db %s for %s err: %v", instanceName, r, err)
			continue
		}
		group.Reads = append(group.Reads, rdb)
	}

	if len(readHosts) > 0 && len(group.Reads) == 0 {
		return fmt.Errorf("no readable db available for %s", instanceName)
	}

	m.mu.Lock()
	m.groups[instanceName] = group
	m.mu.Unlock()
	return nil
}

// GetWriteDB 返回写库
func (m *Manager) GetWriteDB(name string) (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	g, ok := m.groups[name]
	if !ok || g.Write == nil {
		return nil, fmt.Errorf("no write db for %s", name)
	}
	return g.Write, nil
}

// GetReadDB 返回读库（轮询）
func (m *Manager) GetReadDB(name string) (*gorm.DB, error) {
	m.mu.RLock()
	g, ok := m.groups[name]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no db group %s", name)
	}
	g.mu.RLock()
	n := len(g.Reads)
	g.mu.RUnlock()

	if n == 0 {
		if g.Write != nil {
			return g.Write, nil
		}
		return nil, fmt.Errorf("no read/write db available for %s", name)
	}

	idx := int(atomic.AddUint32(&g.rr, 1)) % n
	g.mu.RLock()
	db := g.Reads[idx]
	g.mu.RUnlock()
	return db, nil
}

// CloseAll 优雅关闭
func (m *Manager) CloseAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, g := range m.groups {
		if g.Write != nil {
			if sqlDB, err := g.Write.DB(); err == nil {
				sqlDB.Close()
			}
			g.Write = nil
		}
		for i, r := range g.Reads {
			if sqlDB, err := r.DB(); err == nil {
				sqlDB.Close()
			}
			g.Reads[i] = nil
		}
		g.Reads = nil
	}
}

func (m *Manager) GetDefaultWriteDB() (*gorm.DB, error) {
	var err error
	DefaultWriteDB, err = m.GetWriteDB("default")
	if err != nil {
		return nil, err
	}
	return DefaultWriteDB, nil
}

// StartHealthCheck 启动健康检查
func (m *Manager) StartHealthCheck() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.healthCheck()
		}
	}
}

// healthCheck 检查数据库连接健康状态（仅监控）
func (m *Manager) healthCheck() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	//hlog.Info("[mysql healthcheck] start")

	for name, group := range m.groups {
		// 检查写库
		if group.Write != nil {
			if sqlDB, err := group.Write.DB(); err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := sqlDB.PingContext(ctx); err != nil {
					hlog.Warnf("[healthcheck] write db %s ping failed: %v (will auto-reconnect on next query)", name, err)
				}
				cancel()
			}
		}

		// 检查读库
		for i, readDB := range group.Reads {
			if readDB != nil {
				if sqlDB, err := readDB.DB(); err == nil {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					if err := sqlDB.PingContext(ctx); err != nil {
						hlog.Warnf("[healthcheck] read db %s[%d] ping failed: %v (will auto-reconnect on next query)", name, i, err)
					}
					cancel()
				}
			}
		}
	}
}
