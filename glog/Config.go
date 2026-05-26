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
	*zap.Logger
	sugar *zap.SugaredLogger
}

func NewGLogger(log zap.Logger) *GLogger {
	return &GLogger{
		Logger: &log,
		sugar:  log.Sugar(),
	}
}

func (l *GLogger) InfoF(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}
func (l *GLogger) ErrorF(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}
func (l *GLogger) DebugF(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}
func (l *GLogger) WarnF(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}
