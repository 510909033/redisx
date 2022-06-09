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
	fn := hookFn(ctx)

	cmd := rc.Cluster.HSet(ctx, key, value)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//检查给定域 field 是否存在于哈希表 hash 当中。
func (rc RedisCluster) HExists(ctx *gin.Context, key, feild string) (bool, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.HExists(ctx, key, feild)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//HGET 命令在默认情况下返回给定域的值。
//
//如果给定域不存在于哈希表中， 又或者给定的哈希表并不存在， 那么命令返回 nil 。
func (rc RedisCluster) HGet(ctx *gin.Context, key, feild string) (string, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.HGet(ctx, key, feild)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略
//
//返回值：被成功移除的域的数量，不包括被忽略的域。
func (rc RedisCluster) HDel(ctx *gin.Context, key string, feilds ...string) (int64, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.HDel(ctx, key, feilds...)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//哈希表中域的数量。
//
//当 key 不存在时，返回 0 。
func (rc RedisCluster) HLen(ctx *gin.Context, key string) (int64, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.HLen(ctx, key)

	fn(cmd)

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
	fn := hookFn(ctx)

	cmd := rc.Cluster.HIncrBy(ctx, key, field, incr)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

func (rc RedisCluster) HKeys(ctx *gin.Context, key string) ([]string, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.HKeys(ctx, key)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//以列表形式返回哈希表的域和域的值。
//
//当 key 不存在时，返回一个空表。
func (rc RedisCluster) HGetAll(ctx *gin.Context, key string) (map[string]string, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.HGetAll(ctx, key)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}
