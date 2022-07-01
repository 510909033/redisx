package redisx

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"testing"
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

func TestRedisCluster_ClientGetName(t *testing.T) {
	//cluster := getTestClient()
	//ctx := getTestCtx()

	//cluster.ClientGetName(ctx)

}
