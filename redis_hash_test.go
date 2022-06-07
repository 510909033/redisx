package redisx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisCluster_HSet(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err := cluster.HSet(ctx, key, value)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(value)), cnt)

	value = map[string]string{
		"1": "11",
		"2": "22",
		"3": "3333",
		"4": "444",
	}
	cnt, err = cluster.HSet(ctx, key, value)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), cnt)

}

func TestRedisCluster_HExists(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err := cluster.HSet(ctx, key, value)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(len(value)), cnt)

	exists, err := cluster.HExists(ctx, key, "2")
	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = cluster.HExists(ctx, key, "not exist")
	assert.Nil(t, err)
	assert.False(t, exists)

}

func TestRedisCluster_HGet(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err := cluster.HSet(ctx, key, value)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(len(value)), cnt)

	val, err := cluster.HGet(ctx, key, "2")
	assert.Nil(t, err)
	assert.Equal(t, "22", val)

	val, err = cluster.HGet(ctx, key, "not exist")
	assert.Nil(t, err)
	assert.Equal(t, "", val)

	val, err = cluster.HGet(ctx, "key_exists", "not exist")
	assert.Nil(t, err)
	assert.Equal(t, "", val)
}

func TestRedisCluster_HDel(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err := cluster.HSet(ctx, key, value)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(len(value)), cnt)

	val, err := cluster.HDel(ctx, key, "2", "3", "no_key")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), val)

}

func TestRedisCluster_HLen(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	cnt, err := cluster.HLen(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), cnt)

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err = cluster.HSet(ctx, key, value)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(len(value)), cnt)

	val, err := cluster.HLen(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(value)), val)
}

func TestRedisCluster_HIncrBy(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	cnt, err := cluster.HIncrBy(ctx, key, "num", 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(10), cnt)

	_, err = cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err = cluster.HSet(ctx, key, value)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(len(value)), cnt)

	val, err := cluster.HIncrBy(ctx, key, "2", 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(32), val)
}

func TestRedisCluster_HKeys(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	keys, err := cluster.HKeys(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, (0), len(keys))

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err := cluster.HSet(ctx, key, value)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(len(value)), cnt)

	keys, err = cluster.HKeys(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, (len(value)), len(keys))

	val, err := cluster.HIncrBy(ctx, key, "new_key", 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(10), val)

	keys, err = cluster.HKeys(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, (len(value) + 1), len(keys))
}

func TestRedisCluster_HGetAll(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, nil, err, err)

	redisValues, err := cluster.HGetAll(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{}, redisValues)

	value := map[string]string{
		"1": "11",
		"2": "22",
		"3": "333",
	}
	cnt, err := cluster.HSet(ctx, key, value)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(len(value)), cnt)

	redisValues, err = cluster.HGetAll(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, value, redisValues)

}
