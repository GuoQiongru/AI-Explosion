package main

import (
	"TikTok/controller"
	"TikTok/dao"
	"TikTok/middleware/ftp"
	"TikTok/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	initDeps()
}

func initDeps() {
	//初始化数据库
	dao.Init()
	ftp.InitFTP()

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

}
