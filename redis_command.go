package redisx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var isDebug = true

type Bbtctx = *gin.Context

type client redis.UniversalClient

var logDebug func(ctx Bbtctx, format string, vals ...interface{})
var logErr func(ctx Bbtctx, format string, vals ...interface{})
var errRedisLibrary = errors.New("redis库方法发生未知错误")
var errByCaller = errors.New("使用方调用错误")

func init() {
	logDebug = func(ctx Bbtctx, format string, vals ...interface{}) {
		log.Printf("DEBUG - "+format, vals...)
	}
	logErr = func(ctx Bbtctx, format string, vals ...interface{}) {
		log.Printf("ERROR - "+format, vals...)
	}
}

type RedisCluster struct {
	Cluster client
	Group   string
	options *redis.UniversalOptions
	logger  func(format string, vals ...interface{})
}

// redis key是否存在
func (rc *RedisCluster) Exist(ctx *gin.Context, key string) (bool, error) {
	if !validKey(key) {
		return false, errEmptyKey
	}
	cmd := rc.Cluster.Exists(ctx, key)

	return cmd.Val() == 1, cmd.Err()
}

// keys == 1 时, 0表示key不存在，1表示存在
//
// len(keys) > 1时, int64表示存在的key的数量
func (rc *RedisCluster) Exists(ctx *gin.Context, keys ...string) (int64, error) {

	cmd := rc.Cluster.Exists(ctx, keys...)

	return cmd.Val(), cmd.Err()
}

// Get 获取与key关联的val，第二个返回值表示key是否存在
func (rc *RedisCluster) Get(ctx *gin.Context, key string) (string, bool, error) {
	cmd := rc.Cluster.Get(ctx, key)

	return cmd.Val(), cmd.Err() != redis.Nil, hasErr(cmd)
}

// Delete 删除key，bool表示key是否删除成功，true说明key存在，false表示待删的key不存在
//
// error表示出错
func (rc *RedisCluster) Delete(ctx *gin.Context, key string) (bool, error) {

	cmd := rc.Cluster.Del(ctx, key)

	return cmd.Val() == 1, hasErr(cmd)
}

// 设置key的有效期
//
// seconds 秒数,  如果为0， 会删除这个key。
//
func (rc *RedisCluster) Expire(ctx *gin.Context, key string, seconds int) (bool, error) {

	cmd := rc.Cluster.Expire(ctx, key, time.Second*time.Duration(seconds))

	return cmd.Val(), hasErr(cmd)
}

// 设置key的有效期
//
// timeSet  最终会精确到秒级别(time.Unix())
func (rc *RedisCluster) ExpireAt(ctx *gin.Context, key string, timeSet time.Time) (bool, error) {

	cmd := rc.Cluster.ExpireAt(ctx, key, timeSet)

	return cmd.Val(), hasErr(cmd)
}

//当 key 不存在时，返回 -2 。
//
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
//
// 否则，以秒为单位，返回 key 的剩余生存时间。
func (rc *RedisCluster) TTL(ctx *gin.Context, key string) (time.Duration, error) {

	cmd := rc.Cluster.TTL(ctx, key)

	return cmd.Val(), hasErr(cmd)
}

//为键 key 储存的数字值减去一。
//
//如果键 key 不存在， 那么键 key 的值会先被初始化为 0 ， 然后再执行 DECR 操作。
//
//如果键 key 储存的值不能被解释为数字， 那么 DECR 命令将返回一个错误。
//
//本操作的值限制在 64 位(bit)有符号数字表示之内。
//
//DECR 命令会返回键 key 在执行减一操作之后的值。
func (rc *RedisCluster) Decr(ctx *gin.Context, key string) (int64, error) {

	cmd := rc.Cluster.Decr(ctx, key)

	return cmd.Val(), hasErr(cmd)
}

//将键 key 储存的整数值减去减量 decrement 。
//
//如果键 key 不存在， 那么键 key 的值会先被初始化为 0 ， 然后再执行 DECRBY 命令。
//
//如果键 key 储存的值不能被解释为数字， 那么 DECRBY 命令将返回一个错误。
//
//本操作的值限制在 64 位(bit)有符号数字表示之内。
//
//关于更多递增(increment) / 递减(decrement)操作的更多信息， 请参见 INCR 命令的文档。
//
//返回值
//
//DECRBY 命令会返回键在执行减法操作之后的值。
func (rc *RedisCluster) DecrBy(ctx *gin.Context, key string, decrement int64) (int64, error) {

	cmd := rc.Cluster.DecrBy(ctx, key, decrement)

	return cmd.Val(), hasErr(cmd)
}

//为键 key 储存的值加上浮点数增量 increment 。
//
//如果键 key 不存在， 那么 INCRBYFLOAT 会先将键 key 的值设为 0 ， 然后再执行加法操作。
//
//如果命令执行成功， 那么键 key 的值会被更新为执行加法计算之后的新值， 并且新值会以字符串的形式返回给调用者。
func (rc *RedisCluster) IncrByFloat(ctx *gin.Context, key string, value float64) (float64, error) {

	cmd := rc.Cluster.IncrByFloat(ctx, key, value)

	return cmd.Val(), hasErr(cmd)
}

//
// value不能为空指针
// expiration 0ns表示对应的key，不会设置有效期， 永久有效
//
// expiration 小于0 会报错
//
// error表示出错
func (rc RedisCluster) Set(ctx *gin.Context, key string, value string, expiration time.Duration) error {
	cmd := rc.Cluster.Set(ctx, key, value, expiration)

	return hasErr(cmd)
}

// SetNx 当key不存在时,设置key的值
//
// expiration 0 表示对应的key，不会设置有效期
//
// bool true表示设置成功,即key原来是不存在的， false表示值已经存在，即新设置的值不生效
//
// error表示出错
func (rc RedisCluster) SetNx(ctx *gin.Context, key string, value interface{}, expiration time.Duration) (bool, error) {

	//bytesData, err := json.Marshal(value)
	//if err != nil {
	//	err = fmt.Errorf("GetAny, json.Marshal错误，  key=%s, expiration=%+v, err=%w", key, expiration, err)
	//	logErr(ctx, "%+v", err)
	//	return false, err
	//}

	//cmd := rc.Cluster.SetNX(ctx, key, string(bytesData), expiration)
	cmd := rc.Cluster.SetNX(ctx, key, value, expiration)

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

	cmd := rc.Cluster.IncrBy(ctx, key, value)

	return cmd.Val(), hasErr(cmd)
}
