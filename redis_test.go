package redisx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
	"log"
	"math"
	"strings"
	"testing"
	"time"
)

var addrs = []string{
	"172.20.10.40:36381",
	"172.20.10.40:36382",
	"172.20.10.40:36384",
}

func getTestClient() RedisCluster {
	password := "bitnami"
	log.SetFlags(log.Lshortfile)

	if true {
		str := `172.20.11.141:6379,172.20.11.140:6379,172.20.11.23:6379,172.20.11.26:6379,172.20.11.2:6379,172.20.11.237:6379`
		addrs = strings.Split(str, ",")
		password = ""
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          addrs,
		NewClient:      nil,
		MaxRedirects:   0,
		ReadOnly:       false,
		RouteByLatency: false,
		RouteRandomly:  false,
		ClusterSlots:   nil,
		Dialer:         nil,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			boolCmd := cn.ClientSetName(ctx, "client_name_"+time.Now().Format("2006-01-02__15:04:05"))
			return boolCmd.Err()
		},
		Username:           "",
		Password:           password,
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

	client.AddHook(redisHook{})

	return RedisCluster{
		client: client,
		Group:  "demo_group",
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

func getRedisKey() string {
	return "demo_del"
}

func getRedisLongKey() string {
	return strings.Repeat("demo_del", 500)
}

func deleteKey(key string) {
	cluster := getTestClient()
	ctx := getTestCtx()
	_, err := cluster.Delete(ctx, key)
	if err != nil {
		panic(err)
	}
}

func TestRedisCluster_ConfigGet(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()
	sliceCmd := cluster.client.ConfigGet(ctx, "*")

	var key, val interface{}
	for k, v := range sliceCmd.Val() {
		if k%2 == 0 {
			key = v
			continue
		} else {
			val = v
		}
		t.Logf("k=%s, v=%+v", key, val)
	}

	return
}
func TestRedisCluster_Delete(t *testing.T) {

	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_del"
	_, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	cluster.client.Set(ctx, key, "11", time.Minute)

	cnt, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, true, cnt, err)

	cnt, err = cluster.Delete(ctx, key)
	assert.Equal(t, false, cnt, err)

}

func TestRedisCluster_Exist(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	//??????key
	ret, err := cluster.Exist(ctx, "")
	assert.NotNil(t, err)
	assert.Equal(t, false, ret)

	key := getRedisKey()

	//?????????
	deleteKey(key)
	ret, err = cluster.Exist(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, false, ret)

	//??????
	cluster.client.Set(ctx, key, "11", time.Minute)
	ret, err = cluster.Exist(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, ret)
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

	cluster.client.Set(ctx, key, "11", time.Minute)

	cnt, err = cluster.Exists(ctx, key)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, int64(1), cnt, err)

	cluster.Delete(ctx, key)
	cnt, err = cluster.Exists(ctx, key)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, int64(0), cnt, err)
}

func TestRedisCluster_SetAny(t *testing.T) {

	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"
	_, err := cluster.Delete(ctx, key)
	assert.Equal(t, nil, err, err)

	//????????????
	deleteKey(key)
	err = cluster.Set(ctx, key, "123", time.Minute)
	assert.Equal(t, nil, err, err)

	str, hit, err := cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "123", str)

	//???????????????
	deleteKey(key)
	err = cluster.Set(ctx, key, "str", time.Minute)
	assert.Nil(t, err)

	str, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "str", str)

	// ????????????0
	deleteKey(key)
	err = cluster.Set(ctx, key, "str", 0)
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 500)
	hit, err = cluster.Exist(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)

	duration, err := cluster.TTL(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, time.Duration(-1), duration)

	// ????????????-1
	//todo ??????????????? -1?????????????????????????????????
	//deleteKey(key)
	//err = cluster.Set(ctx, key, "str", -1)
	//assert.NotNil(t, err)

	//duration, err = cluster.TTL(ctx, key)
	//assert.Nil(t, err)
	//assert.Equal(t, time.Duration(-2), duration)
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

	//err = cluster.Set(ctx, key, "123", -1)
	//assert.NotNil(t, err)
	duration, err = cluster.TTL(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, duration, -time.Nanosecond)

	err = cluster.Set(ctx, key, "123", time.Minute)
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

func TestRedisCluster_Todo(t *testing.T) {
	//cluster := getTestClient()
	//ctx := getTestCtx()
	//cluster.Todo(ctx)
}

func TestRedisCluster_Get(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	//??????key
	ret, hit, err := cluster.Get(ctx, "")
	assert.Nil(t, err)
	assert.Equal(t, false, hit)
	assert.Equal(t, "", ret)

	key := getRedisKey()

	//?????????
	deleteKey(key)
	ret, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, false, hit)
	assert.Equal(t, "", ret)

	//??????
	cluster.client.Set(ctx, key, "11", time.Minute)
	ret, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "11", ret)

	//???key
	key = getRedisLongKey()
	deleteKey(key)
	cluster.client.Set(ctx, key, "11", time.Minute)
	ret, hit, err = cluster.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, hit)
	assert.Equal(t, "11", ret)
}

//
//func TestRedisCluster_Set(t *testing.T) {
//
//	cluster := getTestClient()
//	ctx := getTestCtx()
//
//	key := getRedisKey()
//	deleteKey(key)
//
//	//set ??????
//	err := cluster.Set(ctx, key, "", time.Minute)
//	assert.Equal(t, nil, err)
//
//	ret, hit, err := cluster.Get(ctx, key)
//	assert.Nil(t, err)
//	assert.Equal(t, true, hit)
//	assert.Equal(t, "", ret)
//
//	//?????????????????????
//	deleteKey(key)
//	err = cluster.Set(ctx, key, "fadfaf aaaa", time.Minute)
//	assert.Equal(t, nil, err)
//
//	ret, hit, err = cluster.Get(ctx, key)
//	assert.Nil(t, err)
//	assert.Equal(t, true, hit)
//	assert.Equal(t, "fadfaf aaaa", ret)
//
//	//?????? ??????
//	deleteKey(key)
//	err = cluster.Set(ctx, key, 123456, time.Minute)
//	assert.Equal(t, nil, err)
//
//	ret, hit, err = cluster.Get(ctx, key)
//	assert.Nil(t, err)
//	assert.Equal(t, true, hit)
//	assert.Equal(t, "123456", ret)
//
//	//????????????????????? ??????
//	deleteKey(key)
//	timeVal := time.Now()
//	err = cluster.Set(ctx, key, timeVal, time.Minute)
//	assert.Equal(t, nil, err)
//
//	ret, hit, err = cluster.Get(ctx, key)
//	assert.Nil(t, err)
//	assert.Equal(t, true, hit)
//	bytesTime, err := timeVal.MarshalText()
//
//	newT := &time.Time{}
//	err = newT.UnmarshalBinary([]byte(ret))
//	if err != nil {
//		panic(err)
//	}
//	assert.Equal(t, string(bytesTime), ret)
//	assert.Equal(t, timeVal, newT)
//
//	//????????????
//	//deleteKey(key)
//	//newVal := cluster
//	//err = cluster.Set(ctx, key, newVal, time.Minute)
//	//assert.Equal(t, nil, err)
//
//	//ret, hit, err = cluster.Get(ctx, key)
//	//assert.Nil(t, err)
//	//assert.Equal(t, true, hit)
//	//bytesTime, err := timeVal.MarshalText()
//	//
//	//assert.Equal(t, string(bytesTime), ret)
//	//
//
//}
