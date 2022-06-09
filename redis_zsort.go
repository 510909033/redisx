package redisx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Z struct {
	Score  float64
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
func (rc RedisCluster) ZCount(ctx *gin.Context, key string, minScore, maxScore string) (int64, error) {
	fn := hookFn(ctx)

	//strconv.FormatFloat()
	cmd := rc.Cluster.ZCount(ctx, key, minScore, maxScore)

	fn(cmd)

	return cmd.Val(), hasErr(cmd)
}
