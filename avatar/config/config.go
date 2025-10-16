package config

import (
	"avatar/common"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var Cfg *viper.Viper

// DBPoolConfig 池配置（minutes 表示）
//type DBPoolConfig struct {
//	MaxOpenConns    int `mapstructure:"max_open_conns"`
//	MaxIdleConns    int `mapstructure:"max_idle_conns"`
//	ConnMaxLifetime int `mapstructure:"conn_max_lifetime_min"` // minutes
//}
//
//type Hosts struct {
//	Write string   `mapstructure:"write"`
//	Reads []string `mapstructure:"reads"`
//}
//
//type DBParams struct {
//	ParseTime string `mapstructure:"parseTime"`
//	Charset   string `mapstructure:"charset"`
//	Loc       string `mapstructure:"loc"`
//}
//
//type DBCfg struct {
//	Driver         string       `mapstructure:"driver"`
//	DSN            string       `mapstructure:"dsn"`
//	Hosts          Hosts        `mapstructure:"hosts"`
//	User           string       `mapstructure:"user"`
//	Password       string       `mapstructure:"password"`
//	DBName         string       `mapstructure:"dbname"`
//	Params         DBParams     `mapstructure:"params"`
//	Pool           DBPoolConfig `mapstructure:"pool"`
//	ReadWriteSplit bool         `mapstructure:"read_write_split"`
//	Health         struct {
//		PingTimeoutMs int `mapstructure:"ping_timeout_ms"`
//		PingRetries   int `mapstructure:"ping_retries"`
//	} `mapstructure:"health"`
//}
//
//type AppCfg struct {
//	Env      string `mapstructure:"env"`
//	LogLevel string `mapstructure:"log_level"`
//	Port     int    `mapstructure:"port"`
//	LogFile  string `mapstructure:"log_file"`
//}
//
//type Config struct {
//	App      AppCfg           `mapstructure:"app"`
//	Database map[string]DBCfg `mapstructure:"database"`
//}

// LoadConfig 从指定路径加载配置文件，并启用环境变量覆盖
//func LoadConfig(path string) (*Config, error) {
//	v := viper.New()
//	v.SetConfigName(string(common.Env)) // 不要带后缀
//	v.SetConfigType("yaml")
//	if path != "" {
//		// conf 搜索路径
//		v.AddConfigPath(path)
//		// 直接指定配置文件的完整路径（包括文件名和扩展名）
//		// v.SetConfigFile(path)
//		if err := v.ReadInConfig(); err != nil {
//			return nil, err
//		}
//	}
//	// 环境变量前缀（可选）
//	//v.SetEnvPrefix("MYAPP")
//	v.AutomaticEnv() // 自动读取环境变量，如 APP_LOG_LEVEL -> app.log_level
//	Cfg = v
//
//	var cfg Config
//	if err := v.Unmarshal(&cfg); err != nil {
//		return nil, err
//	}
//
//	// 填充默认池参数
//	for name, db := range cfg.Database {
//		if db.Pool.MaxOpenConns == 0 {
//			db.Pool.MaxOpenConns = 20
//		}
//		if db.Pool.MaxIdleConns == 0 {
//			db.Pool.MaxIdleConns = 5
//		}
//		if db.Pool.ConnMaxLifetime == 0 {
//			db.Pool.ConnMaxLifetime = 30
//		}
//		// 默认 health values
//		if db.Health.PingTimeoutMs == 0 {
//			db.Health.PingTimeoutMs = 500
//		}
//		if db.Health.PingRetries == 0 {
//			db.Health.PingRetries = 1
//		}
//		cfg.Database[name] = db
//	}
//
//	// small: add convenient method to convert minute fields if needed elsewhere
//	_ = time.Minute
//
//	return &cfg, nil
//}

func LoadConfig(path string) error {
	v := viper.New()
	v.SetConfigName(string(common.Env)) // 不要带后缀
	v.SetConfigType("yaml")
	if path != "" {
		// conf 搜索路径
		v.AddConfigPath(path)
		// 直接指定配置文件的完整路径（包括文件名和扩展名）
		// v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return err
		}
	}
	v.AddConfigPath(".")
	// 环境变量前缀（可选）
	v.SetEnvPrefix("AVATAR")
	// 把环境变量中的_替换成.，这样可以直接读取环境变量的值
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // 自动读取环境变量，如 APP_LOG_LEVEL -> app.log_level
	Cfg = v
	// small: add convenient method to convert minute fields if needed elsewhere
	_ = time.Minute

	return nil
}
