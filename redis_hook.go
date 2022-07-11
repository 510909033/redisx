package redisx

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var mertic_time_key = "_redis_cmd_consume_time" //命令耗时key
type redisHook struct{}

var _ redis.Hook = (*redisHook)(nil)
var errEmptyKey = errors.New("key的值为空")

//命令执行前
func (redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	bbtctx, ok := ctx.(Bbtctx)
	if !ok {
		if logErr != nil {
			logErr(nil, "ctx.(Bbtctx) 断言失败, name=%s", cmd.Name())
		}
		return ctx, nil
	}
	bbtctx.Set(mertic_time_key, time.Now())
	return ctx, nil
}

//命令执行后
func (redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	bbtctx, ok := ctx.(Bbtctx)
	if !ok {
		if logErr != nil {
			logErr(nil, "ctx.(Bbtctx) 断言失败, name=%s", cmd.Name())
		}
		return nil
	}
	t := bbtctx.GetTime(mertic_time_key)
	if t.IsZero() {
		if logErr != nil {
			logErr(nil, "未获取到命令耗时起始时间, name=%s", cmd.Name())
		}
		return nil
	}

	if isDebug {
		if isRedisNil(cmd) {
			logDebug(bbtctx, "name=%s, 返回了Redis.Nil, args=%+v", cmd.Name(), cmd.Args())
		}
		logDebug(bbtctx, "命令名称=%s, 耗时=%d", cmd.Name(), time.Since(t).Milliseconds())
		//logDebug(bbtctx, "name=%s, cmd=%#v", cmd.Name(), cmd
	}
	//log.Println(cmd.Name(), time.Since(t).Milliseconds())
	//log.Printf("name=%s, cmd=%#v\n", cmd.Name(), cmd)

	return nil
}

func (redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
