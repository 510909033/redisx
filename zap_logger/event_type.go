package zap_logger

const (
	HTTP     = "http"
	MYSQL    = "mysql"
	REDIS    = "redis"
	MEMCACHE = "memcache"
	MONGODB  = "mongodb"
	KAFKA    = "kafka"
	MQ       = "mq"
	ES       = "es"

	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	PANIC = "PANIC"
	FATAL = "FATAL"

	LOG_BIZ_MONITOR = "go_preg_user_biz_logger"
)

// 上报服务明细字段
type LogInfo struct {
	//事件类型
	EventType string
	//耗时，单位秒
	TimeCost float32
	//事件类型
	EventDesc string
	//事件状态
	EventStatus int32
	//目标服务名/主机名/域名
	TargetHost string
	//目标服务处理耗时，单位秒
	ResponseTime float32
	//目标服务返回状态
	ResponseStatus int32
}
