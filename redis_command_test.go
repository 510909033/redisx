package redisx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestRedisCluster_Command(t *testing.T) {
	//cluster := getTestClient()
	//ctx := getTestCtx()
	//cluster.Command(ctx)
	var c1 *gin.Context
	var c2 = context.Background()
	_ = c2
	fn := func(ctx context.Context) {
		log.Println(ctx, ctx == nil, ctx.(*gin.Context) == nil)
	}

	fn(c1)
	//fn(c2)

}

func TestRedisCluster_SetEX(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := getRedisKey()

	deleteKey(key)

	isSuccess, err := cluster.SetEX(ctx, key, "haha", time.Second*3)
	assert.Nil(t, err)
	assert.True(t, isSuccess)

	isSuccess, err = cluster.SetEX(ctx, key, "haha", time.Millisecond)
	assert.Nil(t, err)
	assert.True(t, isSuccess)

}

func TestRedisCluster_DeleteMulti(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	keys := []string{}
	prefix := "{demo_key_}"
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%s,%d", prefix, i)
		keys = append(keys, key)
		deleteKey(key)
		err := cluster.Set(ctx, key, key, time.Minute)
		assert.Nil(t, err)
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%s,%d", prefix, i)
		s, b, err := cluster.Get(ctx, key)
		assert.Nil(t, err)
		assert.True(t, b)
		assert.Equal(t, key, s)
	}

	//Expected nil, but got: "CROSSSLOT Keys in request don't hash to the same slot"
	succNum, err := cluster.DeleteMulti(ctx, keys...)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(keys)), int64(succNum))

}

func TestRedisCluster_DeleteMultiV2(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	keys := []string{}
	prefix := "{demo_key_}"
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%s,%d", prefix, i)
		keys = append(keys, key)
		deleteKey(key)
		err := cluster.Set(ctx, key, key, time.Minute)
		assert.Nil(t, err)
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%s,%d", prefix, i)
		s, b, err := cluster.Get(ctx, key)
		assert.Nil(t, err)
		assert.True(t, b)
		assert.Equal(t, key, s)
	}

	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("%s,%d", prefix, i)
		deleteKey(key)
	}

	//Expected nil, but got: "CROSSSLOT Keys in request don't hash to the same slot"
	succNum, err := cluster.DeleteMulti(ctx, keys...)
	assert.Nil(t, err)
	assert.Equal(t, int64(70), int64(succNum))

}
