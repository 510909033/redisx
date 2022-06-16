package redisx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Z struct {
	Score  float64 //score可以为负数
	Member string
}

type Member Z
type Members []Z

func (members Members) covert() []*redis.Z {
	zList := make([]*redis.Z, len(members))
	for k, v := range members {
		zList[k] = &redis.Z{
			Score:  v.Score,
			Member: v.Member,
		}
	}
	return zList
}

func (member Member) covert() *redis.Z {
	return &redis.Z{
		Score:  member.Score,
		Member: member.Member,
	}
}

//将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
//
//当 key 存在但不是有序集类型时，返回一个错误。
//
//返回值
//被成功添加的新成员的数量，不包括那些被更新的、已经存在的成员。
func (rc RedisCluster) ZAdd(ctx *gin.Context, key string, members Members) (int64, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZAdd(ctx, key, members.covert()...)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//返回有序集 key 中，成员 member 的 score 值。
//
//如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 0， nil 。
func (rc RedisCluster) ZScore(ctx *gin.Context, key string, member string) (float64, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZScore(ctx, key, member)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//返回值
//member 成员的新 score 值
func (rc RedisCluster) ZIncrBy(ctx *gin.Context, key string, increment float64, member string) (float64, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZIncrBy(ctx, key, increment, member)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//当 key 存在且是有序集类型时，返回有序集的基数。 当 key 不存在时，返回 0 。
func (rc RedisCluster) ZCard(ctx *gin.Context, key string) (int64, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZCard(ctx, key)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (rc RedisCluster) ZCount(ctx *gin.Context, key string, minScore, maxScore float64) (int64, error) {
	fn := hookFn(ctx)

	//strconv.FormatFloat()
	cmd := rc.Cluster.ZCount(ctx, key, FloatToString(minScore), FloatToString(maxScore))

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

//http://redisdoc.com/sorted_set/zrange.html
//
//返回有序集 key 中，指定区间内的成员。
//
//其中成员的位置按 score 值递增(从小到大)来排序。
//
//下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
//你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
//
//超出范围的下标并不会引起错误。 比如说，当 start 的值比有序集的最大下标还要大，或是 start > stop 时，
//ZRANGE 命令只是简单地返回一个空列表。 另一方面，假如 stop 参数的值比有序集的最大下标还要大，
//那么 Redis 将 stop 当作最大下标来处理。
//
//返回指定区间内，有序集成员的列表。
func (rc RedisCluster) ZRange(ctx *gin.Context, key string, start, stop int64) ([]string, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZRange(ctx, key, start, stop)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

// 同ZRange， 但是返回结果包括成员和score的所有信息
func (rc RedisCluster) ZRangeWithScores(ctx *gin.Context, key string, start, stop int64) ([]Z, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZRangeWithScores(ctx, key, start, stop)

	fn(cmd)

	return convertZ(cmd.Val()), hasErr(cmd)
}

// 参数含义同ZRange，
//
//成员的位置按 score 值递减(从大到小)来排序。
func (rc RedisCluster) ZRevRange(ctx *gin.Context, key string, start, stop int64) ([]string, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZRevRange(ctx, key, start, stop)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}

// 同ZRevRange， 但是返回结果包括成员和score的所有信息
func (rc RedisCluster) ZRevRangeWithScores(ctx *gin.Context, key string, start, stop int64) ([]Z, error) {
	fn := hookFn(ctx)

	cmd := rc.Cluster.ZRevRangeWithScores(ctx, key, start, stop)

	fn(cmd)

	return convertZ(cmd.Val()), hasErr(cmd)
}
