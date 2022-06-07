package redisx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

var demo = []string{
	`return redis.call('set','foo','bar')`,
	`return redis.call('set',KEYS[1],'bar')`,
}

var hashLock sync.Mutex
var hashMap = map[string]string{}

const (
	//SCRIPT_GET = fmt.Sprintf(`return redis.call('get','%s')`, key)
	SCRIPT_GET = `return redis.call('get', KEYS[1])`
	//SCRIPT_GET = `return redis.call('get', 'demo_key')`
)

/*
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *Cmd
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *Cmd
	ScriptExists(ctx context.Context, hashes ...string) *BoolSliceCmd
	ScriptFlush(ctx context.Context) *StatusCmd
	ScriptKill(ctx context.Context) *StatusCmd
	ScriptLoad(ctx context.Context, script string) *StringCmd
*/
func (rc RedisCluster) Lua(ctx *gin.Context, script string, keys []string, args ...interface{}) (interface{}, error) {

	hashKey, err := func() (string, error) {
		hashLock.Lock()
		defer hashLock.Unlock()

		var hashKey = hashMap[script]

		if hashMap[script] == "" {

			fn := hookFn(ctx)
			loadCmd := rc.Cluster.ScriptLoad(ctx, script)
			fn(loadCmd)

			if loadCmd.Err() != nil {
				return "", loadCmd.Err()
			}
			if loadCmd.Val() == "" {
				return "", fmt.Errorf("返回的hashKey==空，不应该出现的结果, script=%s", script)
			}
			hashKey = loadCmd.Val()
			hashMap[script] = hashKey
		}
		log.Println("GET HASH_KEY, ", hashKey)
		return hashKey, nil
	}()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	fn := hookFn(ctx)
	shaCmd := rc.Cluster.EvalSha(ctx, hashKey, keys, args)
	fn(shaCmd)

	log.Println(hashKey, shaCmd.Val(), shaCmd.Err())
	return shaCmd.Val(), shaCmd.Err()
}

func (rc RedisCluster) LuaDemo(ctx *gin.Context, script string, keys []string, args ...interface{}) (interface{}, error) {

	hashKey, err := func() (string, error) {
		hashLock.Lock()
		defer hashLock.Unlock()

		var hashKey string

		fn := hookFn(ctx)
		loadCmd := rc.Cluster.ScriptLoad(ctx, script)
		fn(loadCmd)

		if loadCmd.Err() != nil {
			return "", loadCmd.Err()
		}
		if loadCmd.Val() == "" {
			return "", fmt.Errorf("返回的hashKey==空，不应该出现的结果, script=%s", script)
		}
		hashKey = loadCmd.Val()
		hashMap[script] = hashKey
		log.Println("GET HASH_KEY, ", hashKey)
		return hashKey, nil
	}()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	fn := hookFn(ctx)
	shaCmd := rc.Cluster.EvalSha(ctx, hashKey, keys, args)
	fn(shaCmd)

	log.Println(hashKey, shaCmd.Val(), shaCmd.Err())
	return shaCmd.Val(), shaCmd.Err()
}
