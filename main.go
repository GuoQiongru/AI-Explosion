package main

import (
	"TikTok/controller"
	"TikTok/dao"

	"TikTok/middleware/ffmpeg"
	"TikTok/middleware/ftp"
	"TikTok/middleware/jwt"
	"TikTok/middleware/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	initDeps()
}

func initDeps() {
	dao.Init()
	ftp.InitFTP()
	ffmpeg.InitSSH()
	redis.InitRedis()
	r := gin.Default()
	initRouter(r)
	r.Run()
}

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	apiRouter.POST("/publish/action/", jwt.AuthWithForm(), controller.Publish)
	apiRouter.GET("/feed/", jwt.SoftAuth(), controller.Feed)
	apiRouter.GET("/publish/list/", jwt.Auth(), controller.PublishList)

	apiRouter.POST("/favorite/action/", jwt.Auth(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", jwt.Auth(), controller.GetFavouriteList)

	apiRouter.POST("/comment/action/", jwt.Auth(), controller.CommentAction)
	apiRouter.GET("/comment/list/", jwt.SoftAuth(), controller.CommentList)

}
