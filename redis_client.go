package redisx

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func (rc *RedisCluster) ClientList(ctx *gin.Context) {

	list := rc.client.ClientList(ctx)
	/*
		<*>client list: id=3 addr=192.168.112.3:41690 laddr=192.168.112.6:6379 fd=13 name= age=261116 idle=1 flags=S db=0 sub=0 psub=0 multi=-1 qbuf=0 qbuf-free=0 argv-mem=0 obl=0 oll=0 omem=0 tot-mem=20512 events=r cmd=replconf user=default redir=-1
		id=6 addr=192.168.112.1:62679 laddr=192.168.112.6:6379 fd=19 name= age=0 idle=0 flags=N db=0 sub=0 psub=0 multi=-1 qbuf=26 qbuf-free=40928 argv-mem=10 obl=0 oll=0 omem=0 tot-mem=61466 events=r cmd=client user=default redir=-1
	*/

	spew.Println(list)

}

func (rc *RedisCluster) ClientID(ctx *gin.Context) {

	list := rc.client.ClientID(ctx)

	spew.Println(list)

}

// 返回 CLIENT SETNAME 命令为连接设置的名字。
//
// 因为新创建的连接默认是没有名字的， 对于没有名字的连接， CLIENT GETNAME 返回空白回复。
//
// OnConnect 可以设置名字
//
//返回值
//如果连接没有设置名字，那么返回空白回复； 如果有设置名字，那么返回名字。
func (rc *RedisCluster) ClientGetName(ctx *gin.Context) {

	/*
		<*>client getname: redis: nil
	*/
	list := rc.client.ClientGetName(ctx)

	spew.Println(list)
}

func (rc *RedisCluster) ClusterSlots(ctx *gin.Context) ([]redis.ClusterSlot, error) {

	/*
			[{5461 10922 [{ed7e032c93c465e9b16a1a3018f9a0842bd77f25 172.20.11.140:6379}
		{3d84e698764bf24441cfae40fc59dca85dbcfafa 172.20.11.237:6379}]}
		{0 5460 [{86ae01e392f3c061805a276f999c7cd26c47696f 172.20.11.23:6379}
		{762a848c4cce0c364b3c61846fae810996b888e0 172.20.11.2:6379}]}
		{10923 16383 [{738b7618a52007774bc1345daac007f4a03583b0 172.20.11.26:6379}
		{f01bb0ef8ea2ffa9a0ae411dfa9716927a93c432 172.20.11.141:6379}]}]
	*/
	list := rc.client.ClusterSlots(ctx)

	logDebug(ctx, "ClusterSlots 切片长度=%d", len(list.Val()))

	spew.Println(list)
	return list.Result()
}

func (rc *RedisCluster) MemoryUsage(ctx *gin.Context, key string, samples ...int) {

	cmd := rc.client.MemoryUsage(ctx, key, samples...)

	spew.Println(cmd)

}

func (rc *RedisCluster) ClusterKeySlot(ctx *gin.Context, key string) (slot int64, err error) {

	cmd := rc.client.ClusterKeySlot(ctx, key)

	return cmd.Result()
}
func (rc *RedisCluster) DebugObject(ctx *gin.Context, key string) (val string, err error) {

	cmd := rc.client.DebugObject(ctx, key)

	return cmd.Result()
}
