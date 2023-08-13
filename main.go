package main

import (
  	"TikTok/dao"
	"TikTok/controller"
	"TikTok/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func main() {
	initDeps()
}

func initDeps() {
  //初始化数据库
  dao.Init()
  r:= gin.Default()
  initRouter(r)
  r.Run()
}

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/video/publish/", controller.Publish)
	apiRouter.GET("/video/feed/", controller.Feed)
	apiRouter.GET("/video/publishList/", controller.PublishList)
}
