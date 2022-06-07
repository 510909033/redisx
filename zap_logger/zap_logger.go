package zap_logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"runtime"
)

type LogDetail struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	BackupNum  int    `yaml:"backup_num"`
	MaxSize    int    `yaml:"max_size"`
	RetainDays int    `yaml:"retain_days"`
}
type LogConf struct {
	Default    LogDetail `yaml:"default"`
	Trace      LogDetail `yaml:"trace"`
	BackupNum  int       `yaml:"backup_num"`
	MaxSize    int       `yaml:"max_size"`
	RetainDays int       `yaml:"retain_days"`
}

const (
	logTmFmtWithMS = "2006-01-02 15:04:05.000"
)

var Logger *zap.Logger
var ZapWriter zapcore.WriteSyncer

// Init 配置日志模块
func init() {

	//defaultLog := &LogDetail{}
	//var level zapcore.Level
	//
	//if level.UnmarshalText([]byte(defaultLog.Level)) != nil {
	//	level = zapcore.InfoLevel
	//}
	// 自定义时间输出格式
	//customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString("[" + t.Format(logTmFmtWithMS) + "]")
	//}
	// 自定义日志级别显示
	//customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString("[" + level.CapitalString() + "]")
	//}
	// 自定义文件：行号输出项
	//customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString("[" + caller.FullPath() + "]")
	//}
	//encoderConfig := zapcore.EncoderConfig{
	//	//CallerKey:      "file", // 打印文件名和行数
	//	//LevelKey:       "level",
	//	//NameKey:        "name",
	//	//TimeKey:        "time",
	//	MessageKey: "msg",
	//	//StacktraceKey:  "stack",
	//	LineEnding: zapcore.DefaultLineEnding,
	//	//EncodeLevel:    customLevelEncoder,
	//	//EncodeTime:     customTimeEncoder,
	//	EncodeDuration: zapcore.StringDurationEncoder,
	//	//EncodeCaller:   customCallerEncoder,
	//}
	//ZapWriter = zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   defaultLog.File,
	//	MaxSize:    defaultLog.MaxSize, // megabytes
	//	MaxBackups: defaultLog.BackupNum,
	//	MaxAge:     defaultLog.RetainDays, // days
	//	LocalTime:  true,
	//})
	//core := zapcore.NewTee(zapcore.NewCore(
	//	zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(ZapWriter), level))
	//zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(ZapWriter), level))
	//Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger, _ = zap.NewDevelopment(zap.AddCallerSkip(2))
}

func Debug(ctx *gin.Context, category, msg string) {
	Logger.Debug(fmt.Sprintf("%s %s", category, msg))
}
func Debugf(ctx *gin.Context, category, msg string, value ...interface{}) {
	Debug(ctx, category, fmt.Sprintf(msg, value...))
}

func DebugUnify(ctx *gin.Context, category string, logInfo *LogInfo) {
	logInfo.EventDesc = fmt.Sprintf("%s %s", category, logInfo.EventDesc)
	Logger.Debug(LogUnify(ctx, DEBUG, logInfo))
}

func Info(ctx *gin.Context, category, msg string) {
	Logger.Info(fmt.Sprintf("%s %s", category, msg))
}

func Infof(ctx *gin.Context, category, msg string, value ...interface{}) {
	Info(ctx, category, fmt.Sprintf(msg, value...))
}

func InfoUnify(ctx *gin.Context, category string, logInfo *LogInfo) {
	logInfo.EventDesc = fmt.Sprintf("%s %s", category, logInfo.EventDesc)
	Logger.Info(LogUnify(ctx, INFO, logInfo))
}

func Warn(ctx *gin.Context, category, msg string) {
	Logger.Warn(fmt.Sprintf("%s %s", category, msg))
}

func Warnf(ctx *gin.Context, category, msg string, value ...interface{}) {
	Warn(ctx, category, fmt.Sprintf(msg, value...))
}

func WarnUnify(ctx *gin.Context, category string, logInfo *LogInfo) {
	logInfo.EventDesc = fmt.Sprintf("%s %s", category, logInfo.EventDesc)
	Logger.Warn(LogUnify(ctx, WARN, logInfo))
}

func Error(ctx *gin.Context, category, msg string) {
	Logger.Error(fmt.Sprintf("%s %s", category, msg))
}

func Errorf(ctx *gin.Context, category, msg string, value ...interface{}) {
	Error(ctx, category, fmt.Sprintf(msg, value...))
}

func ErrorUnify(ctx *gin.Context, category string, logInfo *LogInfo) {
	logInfo.EventDesc = fmt.Sprintf("%s %s", category, logInfo.EventDesc)
	Logger.Error(LogUnify(ctx, ERROR, logInfo))
}

func Panic(ctx *gin.Context, category, msg string) {
	Logger.Panic(fmt.Sprintf("%s %s", category, msg))
}

func Panicf(ctx *gin.Context, category, msg string, value ...interface{}) {
	Panic(ctx, category, fmt.Sprintf(msg, value...))
}

func PanicUnify(ctx *gin.Context, category string, logInfo *LogInfo) {
	logInfo.EventDesc = fmt.Sprintf("%s %s", category, logInfo.EventDesc)
	Logger.Panic(LogUnify(ctx, PANIC, logInfo))
}

func Fatal(ctx *gin.Context, category, msg string) {
	Logger.Fatal(fmt.Sprintf("%s %s", category, msg))
}

//LogUnify进行日志统一化
func LogUnify(ctx *gin.Context, levelStr string, logInfo *LogInfo) string {

	//构造日志字段
	//dataTime := time.Now().Format(bgf_util.TIME_FORMAT_YMDHIS_LAYOUT)
	dataTime := " "
	//获取appName,requestId,clientIp
	var appName, requestId, clientIp, serverIp string
	if ctx != nil {
		//serverIp = middleware.GetServerIp(ctx)
		//appName = config.AppConfig.AppName
		//requestId = middleware.GetRequestId(ctx)
		//clientIp = middleware.GetClientIp(ctx)
	}
	//运行时变量
	var file = "???"
	var callerName = "???"
	var line int
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(2)
	if ok {
		file = filepath.Base(file)
		callerName = runtime.FuncForPC(pc).Name()
	}
	//调用方传递的信息 http://space.babytree-inc.com/pages/viewpage.action?pageId=28521324
	eventType := logInfo.EventType
	timeCost := logInfo.TimeCost
	eventDesc := logInfo.EventDesc
	eventStatus := logInfo.EventStatus
	targetHost := logInfo.TargetHost
	responseTime := logInfo.ResponseTime
	responseStatus := logInfo.ResponseStatus
	//打印日志
	return fmt.Sprintf("%s|%s|%s|%s|%s:%d|%s|%s|%s|%v|%v|%v|%v|%v|%v|%v|%v", dataTime, levelStr, serverIp, appName, file, line, clientIp, clientIp, callerName,
		timeCost, eventType, eventStatus, targetHost, responseTime, responseStatus, requestId, eventDesc)
}
