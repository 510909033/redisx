package redisx

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestRedisCluster_Lua(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"

	cluster.SetString(ctx, key, "some val", time.Second*10)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//_, err := cluster.Lua(ctx, SCRIPT_GET, []string{key})
			_, err := cluster.LuaDemo(ctx, SCRIPT_GET, []string{key})
			assert.Nil(t, err)
			assert.Equal(t, nil, err, err)
		}()
	}
	wg.Wait()

}

func TestRedisCluster_LuaDemo(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := "demo_key"

	cluster.SetString(ctx, key, "some val", time.Second*10)

	for i := 0; i < 10000; i++ {
		//_, err := cluster.Lua(ctx, SCRIPT_GET, []string{key})
		_, err := cluster.LuaDemo(ctx, SCRIPT_GET, []string{key})
		assert.Nil(t, err)
		assert.Equal(t, nil, err, err)
		time.Sleep(time.Second)
	}

}
