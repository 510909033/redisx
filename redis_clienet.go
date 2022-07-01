package redisx

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/http"
	"log"
	"strings"
)

func GetTestRedisClient() RedisCluster {
	password := "bitnami"
	log.SetFlags(log.Lshortfile)
	addrs := []string{}
	if true {
		str := `172.20.11.141:6379,172.20.11.140:6379,172.20.11.23:6379,172.20.11.26:6379,172.20.11.2:6379,172.20.11.237:6379`
		addrs = strings.Split(str, ",")
		password = ""
	}

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
		Cluster: client,
		Group:   "demo_group",
	}
}

func GetTestCtx() *gin.Context {
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
