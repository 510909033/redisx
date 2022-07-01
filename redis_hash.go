package redisx

import (
	"github.com/gin-gonic/gin"
)

//将哈希表 hash 中域 field 的值设置为 value 。
//
//如果给定的哈希表并不存在， 那么一个新的哈希表将被创建并执行 HSET 操作。
//
//如果域 field 已经存在于哈希表中， 那么它的旧值将被新值 value 覆盖。
//
//----
//
//返回值， 只需判断 error即可
//
// ret含义
// 使用新值覆盖了它的旧值 ret不增加
//
// 在哈希表中新创建 field 域并成功为它设置值时， ret加一
func (rc RedisCluster) HSet(ctx *gin.Context, key string, value map[string]string) (int64, error) {

	cmd := rc.client.HSet(ctx, key, value)

	return cmd.Val(), hasErr(cmd)
}

//检查给定域 field 是否存在于哈希表 hash 当中。
func (rc RedisCluster) HExists(ctx *gin.Context, key, feild string) (bool, error) {

	cmd := rc.client.HExists(ctx, key, feild)

	return cmd.Val(), hasErr(cmd)
}

//HGET 命令在默认情况下返回给定域的值。
//
//如果给定域不存在于哈希表中， 又或者给定的哈希表并不存在， 那么命令返回 nil 。
func (rc RedisCluster) HGet(ctx *gin.Context, key, feild string) (string, bool, error) {

	cmd := rc.client.HGet(ctx, key, feild)

	return cmd.Val(), !isRedisNil(cmd), hasErr(cmd)
}

//删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略
//
//返回值：被成功移除的域的数量，不包括被忽略的域。
func (rc RedisCluster) HDel(ctx *gin.Context, key string, feilds ...string) (int64, error) {

	cmd := rc.client.HDel(ctx, key, feilds...)

	return cmd.Val(), hasErr(cmd)
}

//哈希表中域的数量。
//
//当 key 不存在时，返回 0 。
func (rc RedisCluster) HLen(ctx *gin.Context, key string) (int64, error) {

	cmd := rc.client.HLen(ctx, key)

	return cmd.Val(), hasErr(cmd)
}

//为哈希表 key 中的域 field 的值加上增量 increment 。
//
//增量也可以为负数，相当于对给定域进行减法操作。
//
//如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
//
//如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
//
//对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
//
//本操作的值被限制在 64 位(bit)有符号数字表示之内。
//
//返回值：
//执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。
func (rc RedisCluster) HIncrBy(ctx *gin.Context, key, field string, incr int64) (int64, error) {

	cmd := rc.client.HIncrBy(ctx, key, field, incr)

	return cmd.Val(), hasErr(cmd)
}

//为哈希表 key 中的域 field 加上浮点数增量 increment 。
//
//如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
//
//如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
//
//当以下任意一个条件发生时，返回一个错误：
//
//	域 field 的值不是字符串类型(因为 redis 中的数字和浮点数都以字符串的形式保存，所以它们都属于字符串类型）
//
//	域 field 当前的值或给定的增量 increment 不能解释(parse)为双精度浮点数(double precision floating point number)
//
//返回值：
//执行加法操作之后 field 域的值。
func (rc RedisCluster) HIncrByFloat(ctx *gin.Context, key, field string, incr float64) (float64, error) {

	cmd := rc.client.HIncrByFloat(ctx, key, field, incr)

	return cmd.Val(), hasErr(cmd)
}

//返回哈希表 key 中的所有域。
func (rc RedisCluster) HKeys(ctx *gin.Context, key string) ([]string, error) {

	cmd := rc.client.HKeys(ctx, key)

	return cmd.Val(), hasErr(cmd)
}

//以列表形式返回哈希表的域和域的值。
//
//当 key 不存在时，返回一个空表。
func (rc RedisCluster) HGetAll(ctx *gin.Context, key string) (map[string]string, error) {

	cmd := rc.client.HGetAll(ctx, key)

	return cmd.Val(), hasErr(cmd)
}

//返回哈希表 key 中，一个或多个给定域的值。
//
//如果给定的域不存在于哈希表，那么返回一个 nil 值。
//
//因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (rc RedisCluster) HMGet(ctx *gin.Context, key string, feilds ...string) ([]interface{}, error) {

	cmd := rc.client.HMGet(ctx, key, feilds...)

	return cmd.Val(), hasErr(cmd)
}

//同时将多个 field-value (域-值)对设置到哈希表 key 中。
//
//此命令会覆盖哈希表中已存在的域。
//
//如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
//
//如果命令执行成功，返回 OK 。
//
//当 key 不是哈希表(hash)类型时，返回一个错误。
func (rc RedisCluster) HMSet(ctx *gin.Context, key string, value ...interface{}) (bool, error) {

	cmd := rc.client.HMSet(ctx, key, value...)

	return cmd.Val(), hasErr(cmd)
}

//当且仅当域 field 尚未存在于哈希表的情况下， 将它的值设置为 value 。
//
//如果给定域已经存在于哈希表当中， 那么命令将放弃执行设置操作。
//
//如果哈希表 hash 不存在， 那么一个新的哈希表将被创建并执行 HSETNX 命令。
//
//返回值
//
//HSETNX 命令在设置成功时返回 1 ， 在给定域已经存在而放弃执行设置操作时返回 0 。
func (rc RedisCluster) HSetNX(ctx *gin.Context, key, field, value string) (setSuccess bool, err error) {

	cmd := rc.client.HSetNX(ctx, key, field, value)

	return cmd.Val(), hasErr(cmd)
}
