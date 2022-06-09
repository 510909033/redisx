package redisx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type client redis.UniversalClient
type RedisCluster struct {
	Cluster client
	Group   string
}

func debug(cmder redis.Cmder) map[string]interface{} {
	return map[string]interface{}{
		"Err==redis.nil": cmder.Err() == redis.Nil,
		"Err":            cmder.Err(),
		"Name":           cmder.Name(),
		//"String":         cmder.String(),
		//"Args":           cmder.Args(),
		"FullName": cmder.FullName(),
	}
}

// 当cmder返回redis.Nil时，实际并不是发生了错误，根据不同的方法有不同的含义
func hasErr(cmder redis.Cmder) error {
	if cmder.Err() != nil && cmder.Err() != redis.Nil {
		return cmder.Err()
	}
	return nil
}

func hookFn(ctx *gin.Context) func(cmder redis.Cmder) {
	st := time.Now()

	return func(cmder redis.Cmder) {

		//intCmd, ok := cmder.(*redis.IntCmd)
		//if ok {
		//	zap_logger.Debugf(ctx, "demo_redis", "intCmd.Val=%+v", intCmd.Val())
		//}

		if cmder.Err() == nil || cmder.Err() == redis.Nil {

		}

		data := debug(cmder)
		//耗时
		data["耗时"] = time.Since(st).Milliseconds() //毫秒
		//todo
		if cmder.Err() != nil {

		}
		//zap_logger.Debugf(ctx, "demo_redis", "data=%+v", data)
	}
}

// keys == 1 时, 0表示key不存在，1表示存在
//
// len(keys) > 1时, int64表示存在的key的数量
func (rc *RedisCluster) Exists(ctx *gin.Context, keys ...string) (int64, error) {
	fn := hookFn(ctx)
	cmd := rc.Cluster.Exists(ctx, keys...)
	fn(cmd)

	return cmd.Val(), cmd.Err()

}

// Get 获取与key关联的val，第二个返回值表示key是否存在
func (rc *RedisCluster) Get(ctx *gin.Context, key string) (string, bool, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.Get(ctx, key)

	fn(cmd)

	return cmd.Val(), cmd.Err() != redis.Nil, hasErr(cmd)
}

// Delete 删除key，bool表示key是否删除成功，true说明key存在，false表示待删的key不存在
//
// error表示出错
func (rc *RedisCluster) Delete(ctx *gin.Context, key string) (bool, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.Del(ctx, key)

	fn(cmd)

	return cmd.Val() == 1, hasErr(cmd)
}

// 设置key的有效期
//
// seconds 秒数,  如果为0， 会删除这个key。
//
func (rc *RedisCluster) Expire(ctx *gin.Context, key string, seconds int) (bool, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.Expire(ctx, key, time.Second*time.Duration(seconds))

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

// 设置key的有效期
//
// timestamp  10位的时间戳，秒级别
func (rc *RedisCluster) ExpireAt(ctx *gin.Context, key string, timestamp int64) (bool, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ExpireAt(ctx, key, time.Unix(timestamp, 0))

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//当 key 不存在时，返回 -2 。
//
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
//
// 否则，以秒为单位，返回 key 的剩余生存时间。
func (rc *RedisCluster) TTL(ctx *gin.Context, key string) (time.Duration, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.TTL(ctx, key)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}
