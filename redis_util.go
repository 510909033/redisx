package redisx

import (
	"github.com/go-redis/redis/v8"
	"strconv"
)

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
