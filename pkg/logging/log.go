package logging

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gin/pkg/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var defaultLogger Logger

type Logger struct {
	*zap.Logger
}

func GetLogger() Logger {
	return defaultLogger
}

// Setup 初始化配置日志模块
func Setup() {
	var level zapcore.Level
	if level.UnmarshalText([]byte(setting.AppSetting.LogLevel)) != nil {
		level = zapcore.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		LevelKey: "level", NameKey: "name", TimeKey: "time",
		MessageKey: "msg", StacktraceKey: "stack", CallerKey: "location",
		LineEnding: zapcore.DefaultLineEnding, EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.StringDurationEncoder, EncodeCaller: zapcore.FullCallerEncoder,
	}
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   setting.AppSetting.LogSavePath,
		MaxSize:    1024,
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
	})
	var cores = []zapcore.Core{zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), level)}
	if setting.AppSetting.LogConsole {
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
	}
	defaultLogger = Logger{zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))}
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Logger.Debug(fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	defaultLogger.Logger.Info(fmt.Sprintf(format, args...))
}

// Json debug使用
func Json(format string, args ...interface{}) {
	var res []interface{}
	for _, arg := range args {
		byteStr, err := json.MarshalIndent(arg, "", "    ")
		if err != nil {
			Error("%v", err)
		}
		res = append(res, string(byteStr))
	}
	defaultLogger.Logger.Debug(fmt.Sprintf(format, res...))
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Logger.Warn(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	defaultLogger.Logger.Error(fmt.Sprintf(format, args...))
}

func Panic(format string, args ...interface{}) {
	defaultLogger.Logger.Panic(fmt.Sprintf(format, args...))
}

// WithCtx 带请求上下文的Logger，可以记录一些额外信息，比如traceId
func WithCtx(ctx *gin.Context) Logger {
	return Logger{defaultLogger.With(zap.String("traceId", ctx.GetString("traceId")))}
}

func (r Logger) Debug(format string, args ...interface{}) {
	r.Logger.Debug(fmt.Sprintf(format, args...))
}

func (r Logger) Info(format string, args ...interface{}) {
	r.Logger.Info(fmt.Sprintf(format, args...))
}

func (r Logger) Warn(format string, args ...interface{}) {
	r.Logger.Warn(fmt.Sprintf(format, args...))
}

func (r Logger) Error(format string, args ...interface{}) {
	r.Logger.Error(fmt.Sprintf(format, args...))
}

func (r Logger) Panic(format string, args ...interface{}) {
	r.Logger.Panic(fmt.Sprintf(format, args...))
}
