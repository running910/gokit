package logger

import (
	"fmt"
	"os"
	"time"

	logs "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log         *zap.SugaredLogger
	atomicLevel = zap.NewAtomicLevel()
	hLog        *httpLog
)

type httpLog struct{}

func init() {
	log = nil
}

func getLevel(levelName string) zapcore.Level {
	var level zapcore.Level
	switch levelName {
	case "debug":
		level = zapcore.DebugLevel
	case "Debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "Info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "Warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "Error":
		level = zapcore.ErrorLevel
	case "Dpainc":
		level = zapcore.DPanicLevel
	case "dpainc":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "Panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	case "Fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}
	return level
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func LogInit(isDev bool, level string, logFile string, maxSizeMB int, maxBackups int, maxAgeDays int, compress bool) {
	if logFile == "-" {
		// to stdout
		logs.SetFlags(logs.Ldate | logs.Ltime | logs.Lshortfile)
		logs.SetOutput(os.Stdout)
		return
	}

	if _, err := os.Stat(logFile); err != nil {
		if os.IsNotExist(err) {
			os.Create(logFile)
		}
	}

	var hook *lumberjack.Logger

	os.Chmod(logFile, 0666)

	hook = &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   compress,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	atomicLevel.SetLevel(zapcore.Level(getLevel(level)))

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook)), // 打印到文件
		atomicLevel, // 日志级别
	)

	if isDev {
		caller := zap.AddCaller()
		dev := zap.Development()
		log = zap.New(core, caller, zap.AddCallerSkip(1), dev).Sugar()
	} else {
		log = zap.New(core).Sugar()
	}
	log.Info("")
	log.Info("----------------------------------------------------------------------")

	hLog = new(httpLog)
}

func SetLogLevel(level string) {
	atomicLevel.SetLevel(zapcore.Level(getLevel(level)))
}

func DPanic(args ...interface{}) {
	if log == nil {
		logs.Panic(args...)
		return
	}
	log.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	if log == nil {
		logs.Panicf(template, args...)
		return
	}
	log.DPanicf(template, args...)
}

func Debug(args ...interface{}) {
	if log == nil {
		logs.Print(args...)
		return
	}
	log.Debugln(args...)
}
func Debugf(template string, args ...interface{}) {
	if log == nil {
		logs.Printf(template, args...)
		return
	}
	log.Debugf(template, args...)
}

func Error(args ...interface{}) {
	if log == nil {
		logs.Print(args...)
		return
	}
	log.Errorln(args...)
}
func Errorf(template string, args ...interface{}) {
	if log == nil {
		logs.Printf(template, args...)
		return
	}
	log.Errorf(template, args...)
}

func Info(args ...interface{}) {
	if log == nil {
		logs.Print(args...)
		return
	}
	log.Infoln(args...)
}
func Infof(template string, args ...interface{}) {
	if log == nil {
		logs.Printf(template, args...)
		return
	}
	log.Infof(template, args...)
}

func Warn(args ...interface{}) {
	if log == nil {
		logs.Print(args...)
		return
	}
	log.Infoln(args...)
}

func Warnf(template string, args ...interface{}) {
	if log == nil {
		logs.Printf(template, args...)
		return
	}
	log.Warnf(template, args...)
}

func GetHttpLog() *httpLog {
	return hLog
}

func (l *httpLog) Printf(template string, args ...interface{}) {
	if log == nil {
		//fmt.Fprintf(os.Stderr, "log is nil\n")
		fmt.Fprintf(os.Stderr, template+"\n", args...)
		return
	}
	log.Infof(template, args...)
}
