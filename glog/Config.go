package glog

import (
	"go.uber.org/zap"
)

type Config struct {
	File       string      `yaml:"file"`
	Console    bool        `yaml:"Console"`
	Dir        string      `yaml:"dir"`
	MaxSize    int         `yaml:"MaxSize"`
	MaxAge     int         `yaml:"MaxAge"`
	MaxBackups int         `yaml:"MaxBackups"`
	Compress   bool        `yaml:"Compress"`
	LogLevel   LoggerLevel `yaml:"logLevel"`
}
type GLogger struct {
	zap.Logger
}

func NewGLogger(log zap.Logger) *GLogger {
	theLogger := &GLogger{log}
	return theLogger
}
func (r GLogger) InfoF(temp string, args ...interface{}) {
	r.Sugar().Infof(temp, args...)
}
func (r GLogger) ErrorF(temp string, args ...interface{}) {
	r.Sugar().Errorf(temp, args...)
}
func (r GLogger) DebugF(temp string, args ...interface{}) {
	r.Sugar().Debugf(temp, args...)
}
