package redisx

import (
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

//判断一个redis的key是否合法
func validKey(key string) bool {
	if key == "" {
		return false
	}
	return true
}

// cmder.Err() == redis.Nil
func isRedisNil(cmder redis.Cmder) bool {
	return cmder.Err() == redis.Nil
}

//返回内容是否是 OK
func isOK(val string) bool {
	return val == "OK"
}

// 当cmder返回redis.Nil时，实际并不是发生了错误，根据不同的方法有不同的含义
func hasErr(cmder redis.Cmder) error {
	if cmder.Err() != nil && cmder.Err() != redis.Nil {
		return cmder.Err()
	}
	if isDebug {
		if cmder.Err() == redis.Nil {
			log.Println("HIT redis.Nil, name=", cmder.Name())
		}
	}
	return nil
}
func FloatToString(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

func convertZ(redisZ []redis.Z) Members {
	if len(redisZ) == 0 {
		return nil
	}
	members := make(Members, 0, len(redisZ))
	for _, v := range redisZ {
		members = append(members, Z{
			Score:  v.Score,
			Member: v.Member.(string),
		})
	}
	return members
}
