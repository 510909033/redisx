package redisx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

//将一个或多个值 value 插入到列表 key 的表头
//
//如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表头：
//
// 比如说，对空列表 mylist 执行命令 LPUSH mylist a b c ，列表的值将是 c b a ，
//
// 这等同于原子性地执行 LPUSH mylist a 、 LPUSH mylist b 和 LPUSH mylist c 三个命令。
//
//如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。
//
//当 key 存在但不是列表类型时，返回一个错误。
//
//返回值
//执行 LPUSH 命令后，列表的长度
func (rc RedisCluster) LPush(ctx *gin.Context, key string, value ...interface{}) (int64, error) {

	cmd := rc.Cluster.LPush(ctx, key, value...)

	return cmd.Val(), hasErr(cmd)
}

//将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。
//
//当 key 不存在时， LPUSHX 命令什么也不做。
//
//返回值
//LPUSHX 命令执行之后，表的长度。
func (rc RedisCluster) LPushX(ctx *gin.Context, key string,
	value ...interface{}) (int64, error) {

	cmd := rc.Cluster.LPushX(ctx, key, value...)

	return cmd.Val(), hasErr(cmd)
}

//将一个或多个值 value 插入到列表 key 的表尾(最右边)。
//
//如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
//
//如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。
//
//当 key 存在但不是列表类型时，返回一个错误
//
//返回值
//执行 RPUSH 操作后，表的长度
func (rc RedisCluster) RPush(ctx *gin.Context, key string, value ...interface{}) (int64, error) {

	cmd := rc.Cluster.RPush(ctx, key, value...)

	return cmd.Val(), hasErr(cmd)
}

//将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。
//
//当 key 不存在时， RPUSHX 命令什么也不做。
//
//返回值
//RPUSHX 命令执行之后，表的长度。
func (rc RedisCluster) RPushX(ctx *gin.Context, key string, value ...interface{}) (int64, error) {

	cmd := rc.Cluster.RPushX(ctx, key, value...)

	return cmd.Val(), hasErr(cmd)
}

//移除并返回列表 key 的头元素。
//
//返回值
//列表的头元素。 当 key 不存在时, bool == false
func (rc RedisCluster) LPop(ctx *gin.Context, key string) (string, bool, error) {

	cmd := rc.Cluster.LPop(ctx, key)

	return cmd.Val(), cmd.Err() != redis.Nil, hasErr(cmd)
}

//移除并返回列表 key 的尾元素。
//
//返回值
//列表的尾元素。 当 key 不存在时，返回 nil 。
func (rc RedisCluster) RPop(ctx *gin.Context, key string) (string, error) {

	cmd := rc.Cluster.RPop(ctx, key)

	return cmd.Val(), hasErr(cmd)
}

//返回列表 key 的长度。
//
//如果 key 不存在，则 key 被解释为一个空列表，返回 0 .
//
//如果 key 不是列表类型，返回一个错误。
//
//返回值
//列表 key 的长度。
func (rc RedisCluster) LLen(ctx *gin.Context, key string) (int64, error) {

	cmd := rc.Cluster.LLen(ctx, key)

	return cmd.Val(), hasErr(cmd)
}

//对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
//
//举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
//
//下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
//
//你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
//
//当 key 不是列表类型时，返回一个错误。
//
//注意LTRIM命令和编程语言区间函数的区别
//假如你有一个包含一百个元素的列表 list ，对该列表执行 LTRIM list 0 10 ，结果是一个包含11个元素的列表，
//
// 这表明 stop 下标也在 LTRIM 命令的取值范围之内(闭区间)，这和某些语言的区间函数可能不一致，
//
// 比如Ruby的 Range.new 、 Array#slice 和Python的 range() 函数。
//
//超出范围的下标
//
//超出范围的下标值不会引起错误。
//
//如果 start 下标比列表的最大下标 end ( LLEN list 减去 1 )还要大，或者 start > stop ， LTRIM 返回一个空列表(因为 LTRIM 已经将整个列表清空)。
//
//如果 stop 下标比 end 下标还要大，Redis将 stop 的值设置为 end 。
//
//返回值
//命令执行成功时，true 。
func (rc RedisCluster) LTrim(ctx *gin.Context, key string, start, stop int64) (bool, error) {

	cmd := rc.Cluster.LTrim(ctx, key, start, stop)

	return isOK(cmd.Val()), hasErr(cmd)
}

//根据参数 count 的值，移除列表中与参数 value 相等的元素。
//
//count 的值可以是以下几种：
//
//count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。
//
//count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。
//
//count = 0 : 移除表中所有与 value 相等的值。
//
//返回值
//被移除元素的数量。 因为不存在的 key 被视作空表(empty list)，所以当 key 不存在时， LREM 命令总是返回 0 。
func (rc RedisCluster) LRem(ctx *gin.Context, key string, count int64, value interface{}) (int64, error) {

	cmd := rc.Cluster.LRem(ctx, key, count, value)

	return cmd.Val(), hasErr(cmd)
}

func (rc RedisCluster) LRange(ctx *gin.Context, key string, start, stop int64) ([]string, error) {

	cmd := rc.Cluster.LRange(ctx, key, start, stop)

	return cmd.Val(), hasErr(cmd)
}
