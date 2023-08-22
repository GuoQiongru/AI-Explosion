package redis

import (
	"TikTok/config"
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RdbLikeUserId *redis.Client  //key:userId,value:VideoId
var RdbLikeVideoId *redis.Client //key:VideoId,value:userId

func InitRedis() {

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPsw,
		DB:       5, //  点赞视频id信息存入 DB5.
	})

	RdbLikeVideoId = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPsw,
		DB:       6, //  点赞用户id信息存入 DB6.
	})

}
