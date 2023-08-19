package main

import (
	"TikTok/controller"
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"TikTok/middleware/ftp"
	"TikTok/middleware/jwt"
	"TikTok/middleware/rabbitmq"
	"TikTok/middleware/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	initDeps()
}

func initDeps() {
	//初始化数据库
	dao.Init()
	ftp.InitFTP()
	ffmpeg.InitSSH()
	redis.InitRedis()
	// 初始化rabbitMQ。
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitLikeRabbitMQ()
	r := gin.Default()
	initRouter(r)
	r.Run()
}

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	apiRouter.POST("/video/publish/action/", jwt.AuthWithForm(), controller.Publish)
	apiRouter.GET("/video/feed/", jwt.SoftAuth(), controller.Feed)
	apiRouter.GET("/video/publish/list/", jwt.Auth(), controller.PublishList)

	apiRouter.POST("/favorite/action/", jwt.Auth(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", jwt.Auth(), controller.GetFavouriteList)

}
