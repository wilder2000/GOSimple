package config

import (
	"time"
)

type AppConfig struct {
	ReadTimeout               time.Duration
	WriteTimeout              time.Duration
	PoolSize                  int
	MaxThread                 int
	DataSource                DBConfig
	Web                       map[string]string
	StaticDir                 StaticDir
	ReportError               bool
	AccessControlAllowOrigin  bool
	AccessControlAllowHost    string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
	Security                  Security
	ExpireTime                time.Duration
	DefaultAvatar             string
	UserService               bool
	UserNames                 []string
	AppSecret                 AppSecret
}
type AppSecret struct {
	AccessKey string
	SecretKey string
}
type DBConfig struct {
	Dir                string
	Type               string
	Name               string
	DSN                string
	MaxOpenConnections int
	MaxIdleConnections int
}
type Security struct {
	Registration     bool
	MaxTryTimes      int
	ForbidAccessTime float64
	DefaultAdminPWD  string
}
type StaticDir struct {
	RelativePath    string
	AbsoluteFileDir string
}

func (dbc *DBConfig) IsMySQL() bool {
	return dbc.Type == DBTYPE_MYSQL
}
func (dbc *DBConfig) IsSQLLite() bool {
	return dbc.Type == DBTYPE_SQLITE
}

const (
	DBTYPE_MYSQL  = "mysql"
	DBTYPE_SQLITE = "sqlite"
	IconHome      = "/images"
)

var (
	AConfig *AppConfig
)

func AvatorHome() string {
	return AConfig.StaticDir.AbsoluteFileDir + IconHome
}
func LoadConfig() {
	AConfig = &AppConfig{}
	yaml, err := ReadYAML(AppConfigFile(), ConfDir())
	if err != nil {
		panic(err)
	}
	yaml.Sub("application").Unmarshal(AConfig)
}

func LoadUserApplicationConfig[T interface{}]() *T {
	model := new(T)
	yaml, err := ReadYAML(AppConfigFile(), ConfDir())
	if err != nil {
		panic(err)
	}
	yaml.Sub("application").Unmarshal(model)
	return model
}

//读取配置文件

func ReadYAMLConfig[T interface{}](conf string) *T {
	model := new(T)
	yaml, err := ReadYAML(conf, ConfDir())
	if err != nil {
		panic(err)
	}
	yaml.Sub("application").Unmarshal(model)
	return model
}
func init() {
	LoadConfig()
}
