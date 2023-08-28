package redis

import (
	"TikTok/config"
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RdbLikeUserId *redis.Client  //key:userId,value:VideoId
var RdbLikeVideoId *redis.Client //key:VideoId,value:userId

var RdbVCid *redis.Client //redis db11 -- video_id + comment_id
var RdbCVid *redis.Client //redis db12 -- comment_id + video_id

func InitRedis() {

	RdbLikeUserId = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPsw,
		DB:       1, //  点赞视频id信息存入 DB5.
	})

	RdbLikeVideoId = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPsw,
		DB:       2, //  点赞用户id信息存入 DB6.
	})

	RdbVCid = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPsw,
		DB:       3, // lsy 选择将video_id中的评论id s存入 DB11.
	})

	RdbCVid = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPsw,
		DB:       4, // lsy 选择将comment_id对应video_id存入 DB12.
	})

}
