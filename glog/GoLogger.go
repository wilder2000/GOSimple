package glog

import (
	"GOSimple/config"
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
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
	Logger = logger()
}

// Logger 初始化一个日志
func logger() *GLogger {

	file, err := config.ReadYAML("log4g.yaml", config.ConfDir())
	if err != nil {
		fmt.Println("load log config failed.", err)
	}
	file.Unmarshal(&LConfig)

	hook := lumberjack.Logger{
		Filename:   LConfig.File,       //日志文件路径
		MaxSize:    LConfig.MaxSize,    //在进行切割之前，日志文件的最大大小（以MB为单位
		MaxAge:     LConfig.MaxAge,     //保留旧文件的最大天数
		MaxBackups: LConfig.MaxBackups, //保留旧文件的最大个数
		Compress:   LConfig.Compress,   //是否压缩/归档旧文件
	}
	var writerSync zapcore.WriteSyncer
	if LConfig.Console {
		fmt.Println("console enabled")
		writerSync = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&hook))
	} else {
		fmt.Println("console disabled")
		writerSync = zapcore.AddSync(&hook)
	}

	// debug->info->warn->error
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
	return NewGLogger(*zLog)
}
