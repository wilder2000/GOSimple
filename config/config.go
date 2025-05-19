package config

import (
	"time"
)

type AppConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MaxThread    int
}

var (
	AConfig *AppConfig
)

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
