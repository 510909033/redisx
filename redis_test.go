package redisx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
	"math"
	"testing"
	"time"
)

var addrs = []string{
	"172.20.10.40:36381",
	"172.20.10.40:36382",
	"172.20.10.40:36384",
}

func getTestClient() RedisCluster {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              addrs,
		NewClient:          nil,
		MaxRedirects:       0,
		ReadOnly:           false,
		RouteByLatency:     false,
		RouteRandomly:      false,
		ClusterSlots:       nil,
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           "bitnami",
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolFIFO:           false,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	})

	return RedisCluster{
		Cluster: client,
		Group:   "demo_group",
	}
}

func getTestCtx() *gin.Context {
	w := &http.TestResponseWriter{
		StatusCode: 0,
		Output:     "",
	}
	_ = w
	//ctx, engine := gin.CreateTestContext(w)
	//_ = engine
	ctx := &gin.Context{}
	return ctx
}

func TestRedisCluster_Delete(t *testing.T) {

	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_del"
	_, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	cluster.Cluster.Set(ctx, key, "11", time.Minute)

	cnt, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, true, cnt, err)

	cnt, err = cluster.Delete(ctx, key)
	assert.Equal(t, false, cnt, err)

}

func TestRedisCluster_Exists(t *testing.T) {

	cluster := getTestClient()
	ctx := getTestCtx()

	cnt, err := cluster.Exists(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), cnt)

	key := "demo_key"
	_, err = cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	cluster.Cluster.Set(ctx, key, "11", time.Minute)

	cnt, err = cluster.Exists(ctx, key)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, int64(1), cnt, err)

	cluster.Delete(ctx, key)
	cnt, err = cluster.Exists(ctx, key)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, int64(0), cnt, err)
}

func TestRedisCluster_Set(t *testing.T) {

	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	err = cluster.SetInt64(ctx, key, 123, time.Minute)
	assert.Equal(t, nil, err, err)

	str, hit, err := cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "123", str)

	err = cluster.SetString(ctx, key, "str", time.Minute)
	assert.Nil(t, err)

	str, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "str", str)
}

func TestRedisCluster_SetNx(t *testing.T) {

	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	_, err = cluster.SetNx(ctx, key, "123", time.Minute)
	assert.Equal(t, nil, err, err)

	str, hit, err := cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "123", str)

	_, err = cluster.SetNx(ctx, key, "str", time.Minute)
	assert.Nil(t, err)

	str, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "123", str)
}

func TestRedisCluster_IncrBy(t *testing.T) {

	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	val, err := cluster.IncrBy(ctx, key, 20)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, int64(20), val)

	str, hit, err := cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "20", str)

	val, err = cluster.IncrBy(ctx, key, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(21), val)

	str, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "21", str)

	val, err = cluster.IncrBy(ctx, key, -9)
	assert.Nil(t, err)
	assert.Equal(t, int64(12), val)

	str, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "12", str)

	val, err = cluster.IncrBy(ctx, key, math.MaxInt64)
	assert.NotNil(t, err)

	str, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "12", str)
}

func TestRedisCluster_TTL(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	duration, err := cluster.TTL(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, -time.Nanosecond*2, duration)

	_, err = cluster.SetNx(ctx, key, "123", 0)
	assert.Equal(t, nil, err, err)
	duration, err = cluster.TTL(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, duration, -time.Nanosecond)

	err = cluster.SetString(ctx, key, "123", -1)
	assert.Equal(t, nil, err, err)
	duration, err = cluster.TTL(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, duration, -time.Nanosecond)

	err = cluster.SetString(ctx, key, "123", time.Minute)
	assert.Equal(t, nil, err, err)
	duration, err = cluster.TTL(ctx, key)
	assert.Nil(t, err)
	assert.Greater(t, duration, -time.Nanosecond)

	_, err = cluster.Expire(ctx, key, 0)
	assert.Nil(t, err)

	cnt, err := cluster.Exists(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), cnt)

}
