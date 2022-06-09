package redisx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisCluster_ZScore(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)

	score, err := cluster.ZScore(ctx, key, "no_member")
	assert.Nil(t, err)
	assert.Equal(t, float64(0), score)

	members := Members{
		{
			Score:  0,
			Member: "0",
		},
		{
			Score:  1,
			Member: "1",
		},
		{
			Score:  2,
			Member: "2",
		},
	}

	num, err := cluster.ZAdd(ctx, key, members)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(members)), num)

	score, err = cluster.ZScore(ctx, key, "no_member")
	assert.Nil(t, err)
	assert.Equal(t, float64(0), score)

	score, err = cluster.ZScore(ctx, key, "1")
	assert.Nil(t, err)
	assert.Equal(t, float64(1), score)

}

func TestRedisCluster_ZIncrBy(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)

	members := Members{
		{
			Score:  0,
			Member: "0",
		},
		{
			Score:  1,
			Member: "1",
		},
		{
			Score:  2,
			Member: "2",
		},
	}

	num, err := cluster.ZAdd(ctx, key, members)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(members)), num)

	score, err := cluster.ZIncrBy(ctx, key, 20, "no_member")
	assert.Nil(t, err)
	assert.Equal(t, float64(20), score)

	score, err = cluster.ZIncrBy(ctx, key, 10, "1")
	assert.Nil(t, err)
	assert.Equal(t, float64(11), score)

}

func TestRedisCluster_ZCard(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)

	cnt, err := cluster.ZCard(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), cnt)

	members := Members{
		{
			Score:  0,
			Member: "0",
		},
		{
			Score:  1,
			Member: "1",
		},
		{
			Score:  2,
			Member: "2",
		},
	}

	num, err := cluster.ZAdd(ctx, key, members)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(members)), num)

	cnt, err = cluster.ZCard(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(members)), cnt)

}

func TestRedisCluster_ZCount(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)

	members := Members{
		{
			Score:  10.22222,
			Member: "0",
		},
		{
			Score:  21.88888,
			Member: "1",
		},
		{
			Score:  32.55555,
			Member: "2",
		},
	}

	num, err := cluster.ZAdd(ctx, key, members)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(members)), num)

	cnt, err := cluster.ZCount(ctx, key, "10", "25")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), cnt)

	cnt, err = cluster.ZCount(ctx, key, "10.22222", "10.22222")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), cnt)

	cnt, err = cluster.ZCount(ctx, key, "0", "9.9999")
	assert.Nil(t, err)
	assert.Equal(t, int64(0), cnt)

	cnt, err = cluster.ZCount(ctx, key, "9.9999", "10.22222")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), cnt)

}
