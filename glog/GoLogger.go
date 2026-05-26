package glog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/wilder2000/GOSimple/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LevelDebug LoggerLevel = "debug"
const LevelInfo LoggerLevel = "info"
const LevelError LoggerLevel = "error"

type LoggerLevel string

var Logger *GLogger
var LConfig *Config

func init() {
	var err error
	Logger, err = newLogger()
	if err != nil {
		panic(fmt.Sprintf("glog init failed: %v", err))
	}
}

func newLogger() (*GLogger, error) {
	file, err := config.ReadYAML("log4g.yaml", config.ConfDir())
	if err != nil {
		return nil, fmt.Errorf("load log config failed: %w", err)
	}
	LConfig = &Config{}
	if err := file.Unmarshal(LConfig); err != nil {
		return nil, fmt.Errorf("unmarshal log config failed: %w", err)
	}

	logPath := LConfig.File
	if LConfig.Dir != "" {
		logPath = filepath.Join(LConfig.Dir, LConfig.File)
	}

	hook := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    LConfig.MaxSize,
		MaxAge:     LConfig.MaxAge,
		MaxBackups: LConfig.MaxBackups,
		Compress:   LConfig.Compress,
	}

	var writerSync zapcore.WriteSyncer
	if LConfig.Console {
		writerSync = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&hook))
	} else {
		writerSync = zapcore.AddSync(&hook)
	}

	var level zapcore.Level
	switch LConfig.LogLevel {
	case LevelDebug:
		level = zap.DebugLevel
	case LevelInfo:
		level = zap.InfoLevel
	case LevelError:
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		writerSync,
		level,
	)

	zLog := zap.New(core)
	return NewGLogger(*zLog), nil
}
