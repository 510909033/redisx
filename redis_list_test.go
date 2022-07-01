package redisx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var listKey = "demo_list_key"

func TestRedisCluster_LPop(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := listKey

	deleteKey(key)

	//不存在的key
	val, hit, err := cluster.LPop(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, "", val)
	assert.Equal(t, false, hit)

	deleteKey(key)
	cnt, err := cluster.LPushX(ctx, key, 1, 2, 3)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), cnt)

	deleteKey(key)
	cnt, err = cluster.LPush(ctx, key, 1, 2, 3) //实际存储 3 2 1
	assert.Nil(t, err)
	assert.Equal(t, int64(3), cnt)

	val, hit, err = cluster.LPop(ctx, key) // 2 1
	assert.Nil(t, err)
	assert.Equal(t, "3", val)
	assert.Equal(t, true, hit)

	cnt, err = cluster.LPushX(ctx, key, 1, 2, 3) // 3 2 1 2 1
	assert.Nil(t, err)
	assert.Equal(t, int64(5), cnt)

	valList, err := cluster.LRange(ctx, key, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, []string{"3", "2", "1", "2", "1"}, valList)

	//ltrim
	ok, err := cluster.LTrim(ctx, key, 1, 3) //保留3个元素   // 2 1 2
	assert.Nil(t, err)
	assert.Equal(t, true, ok)

	valList, err = cluster.LRange(ctx, key, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, []string{"2", "1", "2"}, valList)

	//lrem
	cluster.RPush(ctx, key, 4, 5, 6, 2)         // 2 1 2 4 5 6 2
	remCnt, err := cluster.LRem(ctx, key, 1, 1) //  2 2 4 5 6 2
	assert.Nil(t, err)
	assert.Equal(t, int64(1), remCnt)

	remCnt, err = cluster.LRem(ctx, key, 0, 0) // 3 2 4 5 6 2
	assert.Nil(t, err)
	assert.Equal(t, int64(0), remCnt)

	//len
	lLen, err := cluster.LLen(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, int64(6), lLen)

	//lrange
	valList, err = cluster.LRange(ctx, key, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, []string{"2", "2", "4", "5", "6", "2"}, valList)

	remCnt, err = cluster.LRem(ctx, key, 0, 2) // 1  4 5 6
	assert.Nil(t, err)
	assert.Equal(t, int64(3), remCnt)

}
