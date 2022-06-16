package redisx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func createZData(t *testing.T) ([]float64, Members) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Nil(t, err)

	data := []float64{
		-1,
		0,
		0.01,
		0.0001,
		1,
		1.001,
		1.9999,
		2,
		2.001,
	}

	var members = make(Members, 0, len(data))
	for _, v := range data {
		members = append(members, Z{
			Score:  v,
			Member: FloatToString(v),
		})
	}

	num, err := cluster.ZAdd(ctx, key, members)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(members)), num)

	mems, err := cluster.ZRangeWithScores(ctx, key, 0, 10000)
	_ = mems
	assert.Nil(t, err)
	//log.Printf("mems=%#v\n", mems)
	//spew.Dump(mems)
	return data, members
}

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

	count, err := cluster.ZCount(ctx, key, 0, 100000)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)

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

	cnt, err := cluster.ZCount(ctx, key, 10, 25)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), cnt)

	cnt, err = cluster.ZCount(ctx, key, 10.22222, 10.22222)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), cnt)

	cnt, err = cluster.ZCount(ctx, key, 0, 9.9999)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), cnt)

	cnt, err = cluster.ZCount(ctx, key, 9.9999, 10.22222)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), cnt)

}

func TestRedisCluster_ZRange(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()
	key := "demo_key"

	scores, members := createZData(t)
	_ = scores

	data, err := cluster.ZRange(ctx, key, 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, members[0].Member, data[0])
	assert.Equal(t, members[1].Member, data[1])

	data, err = cluster.ZRange(ctx, key, 0, 100)
	assert.Nil(t, err)
	assert.Equal(t, len(scores), len(data))
	assert.Equal(t, members[0].Member, data[0])
	assert.Equal(t, members[len(scores)-1].Member, data[len(data)-1])

	data2, err := cluster.ZRange(ctx, key, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, len(scores), len(data2))
	assert.Equal(t, members[0].Member, data[0])
	assert.Equal(t, members[len(scores)-1].Member, data[len(data)-1])
	assert.Equal(t, data, data2)

}
func TestRedisCluster_ZRangeWithScores(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()
	key := "demo_key"

	scores, members := createZData(t)
	_ = scores

	data, err := cluster.ZRevRangeWithScores(ctx, key, 0, 1)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, members[len(members)-1], data[0])
	assert.Equal(t, members[len(members)-2], data[1])

	data, err = cluster.ZRevRangeWithScores(ctx, key, 0, 100)
	assert.Nil(t, err)
	assert.Equal(t, len(scores), len(data))
	assert.Equal(t, members[len(members)-1], data[0])
	assert.Equal(t, members[0], data[len(data)-1])

	data2, err := cluster.ZRevRangeWithScores(ctx, key, 0, -1)
	assert.Nil(t, err)
	assert.Equal(t, len(scores), len(data2))
	assert.Equal(t, members[len(members)-1], data[0])
	assert.Equal(t, members[0], data[len(data)-1])
	assert.Equal(t, data, data2)

}
