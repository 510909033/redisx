package redisx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	crc16V2 "github.com/sigurn/crc16"
)

func TestRedisCluster_ClientList(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	cluster.ClientList(ctx)

}

// 老版本不支持 ？http://redisdoc.com/client_and_server/index.html
func TestRedisCluster_ClientID(t *testing.T) {
	// todo
	//cluster := getTestClient()
	//ctx := getTestCtx()
	//cluster.ClientID(ctx)
}

func TestRedisCluster_ClientGetName(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	cluster.ClientGetName(ctx)
}

func TestRedisCluster_ClusterSlots(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	cluster.ClusterSlots(ctx)
}
func TestRedisCluster_ClusterKeySlot(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := getRedisKey()
	slot, err := cluster.ClusterKeySlot(ctx, key)
	assert.Nil(t, err)

	//crc16Slot := bb_CRC16([]byte(key), 8)
	//
	//assert.Equal(t, slot, int64(crc16Slot%16384))
	//assert.Equal(t, slot, int64(CheckSum([]byte(key))%16384))
	//assert.Equal(t, slot, int64(crc16_v3([]byte(key))%16384))
	//
	//assert.Equal(t, slot, int64(crc16.ChecksumCCITT([]byte(key))%16384))
	//assert.Equal(t, slot, int64(crc16.ChecksumCCITTFalse([]byte(key))%16384))
	//assert.Equal(t, slot, int64(crc16.ChecksumIBM([]byte(key))%16384))
	//assert.Equal(t, slot, int64(crc16.ChecksumMBus([]byte(key))%16384))
	//assert.Equal(t, slot, int64(crc16.ChecksumSCSI([]byte(key))%16384))

	table := crc16V2.MakeTable(crc16V2.CRC16_MAXIM)

	crc := crc16V2.Checksum([]byte(key), table)

	assert.Equal(t, slot, int64(crc%16384))

	fmt.Printf("CRC-16 MAXIM: %X\n", crc)

	// using the standard library hash.Hash interface
	h := crc16V2.New(table)
	h.Write([]byte("Hello world!"))
	fmt.Printf("CRC-16 MAXIM: %X\n", h.Sum16())

}

func TestRedisCluster_DebugObject(t *testing.T) {
	cluster := getTestClient()
	ctx := getTestCtx()

	key := getRedisKey()

	deleteKey(key)

	val, err := cluster.DebugObject(ctx, key)
	assert.Nil(t, err)
	t.Logf("DebugObject, val=%s\n", val)

}
