package redisx

import (
	"github.com/gin-gonic/gin"
	"time"
)

// SetString 设置key的值
//
// expiration 0，-1表示对应的key，不会设置有效期
//
// error表示出错
func (rc RedisCluster) SetString(ctx *gin.Context, key string, value string, expiration time.Duration) error {
	fn := hookFn(ctx)

	cmd := rc.Cluster.Set(ctx, key, value, expiration)

	fn(cmd)

	return hasErr(cmd)
}

// SetInt64 设置key的值
//
// expiration 0，-1表示对应的key，不会设置有效期
//
// error表示出错
func (rc RedisCluster) SetInt64(ctx *gin.Context, key string, value int64, expiration time.Duration) error {
	fn := hookFn(ctx)

	cmd := rc.Cluster.Set(ctx, key, value, expiration)

	fn(cmd)

	return hasErr(cmd)
}

// SetNx 当key不存在时,设置key的值
//
// expiration 0，-1表示对应的key，不会设置有效期
//
// bool true表示设置成功,即key原来是不存在的， false表示值已经存在，即新设置的值不生效
//
// error表示出错
func (rc RedisCluster) SetNx(ctx *gin.Context, key string, value string, expiration time.Duration) (bool, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.SetNX(ctx, key, value, expiration)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//为键 key 储存的数字值加上 value, 负数表示减法操作
//
//如果键 key 不存在， 那么它的值会先被初始化为 0 ， 然后再执行 INCRBY  命令。
//
//如果键 key 储存的值不能被解释为数字， 那么 INCRBY  命令将返回一个错误。
//
//本操作的值限制在 64 位(bit)有符号数字表示之内。
func (rc RedisCluster) IncrBy(ctx *gin.Context, key string, value int64) (int64, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.IncrBy(ctx, key, value)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}
