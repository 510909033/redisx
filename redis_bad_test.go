package redisx

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRedisCluster_Redis_lock(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()
	key := getRedisKey()
	cluster.SetNx(ctx, key, "val", time.Second) // master成功，save未同步后发生了切换

	cluster.SetNx(ctx, key, "val", time.Second) //slave还会成功

}
func TestRedisCluster_Redis_Exception(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()
	var wg sync.WaitGroup
	wg.Add(1)

	useClose := false
	useShutdown := false

	var index int64
	go func() {
		for {
			key := fmt.Sprintf("key_%d", atomic.AddInt64(&index, 1))
			err := cluster.Set(ctx, key, key, time.Second*3)
			if err != nil {
				log.Printf("TEST, set err=%+v", err)
			}

			val, hit, err := cluster.Get(ctx, key)
			if err != nil {
				log.Printf("TEST, get err=%+v", err)
			} else {
				if !hit || val != key {
					log.Printf("TEST, hit=%t, val=%s, key=%s", hit, val, key)
				}
			}

			time.Sleep(time.Millisecond * 300)
		}
	}()

	go func() {
		if !useClose {
			return
		}
		time.Sleep(time.Second * 3)
		log.Println("client.close start")
		err := cluster.client.Close()
		if err != nil {
			log.Printf("client.Close 失败, err=%+v", err)
		}
		log.Println("client.close end")
	}()

	//Shutdown
	go func() {
		if !useShutdown {
			return
		}
		time.Sleep(time.Second * 3)
		log.Println("client.Shutdown start")
		statusCmd := cluster.client.Shutdown(ctx)
		log.Printf("Shutdown, result=%+v", statusCmd.Val())
		if statusCmd.Err() != nil {
			log.Printf("client.Shutdown 失败, err=%+v", statusCmd.Err())
		}

		log.Println("client.Shutdown end")
	}()

	wg.Wait()

	return
}
