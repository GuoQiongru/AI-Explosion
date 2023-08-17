package controller

import (
	"TikTok/dao"
	"TikTok/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []dao.TableVideo `json:"video_list"`
	NextTime  int64            `json:"next_time"`
}

func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	log.Printf("传入的时间" + inputTime)
	var lastTime time.Time
	if inputTime != "" {
		me, _ := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.Unix(me, 0)
	} else {
		lastTime = time.Now()
	}
	log.Printf("获取到时间戳%v", lastTime)
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)

	videoService := GetVideo()
	feed, err := videoService.Feed(lastTime)
	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, userId) 失败：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}

	log.Printf("方法videoService.Feed(lastTime, userId) 成功")
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  time.Now().Unix(),
	})
}

func Publish(c *gin.Context) {
	file, err := c.FormFile("data")
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	fmt.Printf("获取到用户id:%v\n", userId)
	title := c.PostForm("title")
	log.Printf("获取到视频title:%v\n", title)
	if err != nil {
		log.Printf("获取视频流失败:%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoService := GetVideo()
	err = videoService.Publish(c, file, userId, title)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到当前用户id:%v\n", curId)

	videoService := GetVideo()
	feed, err := videoService.List(userId)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  time.Now().Unix(),
	})
}

func GetVideo() service.VideoServiceImpl {
	var videoService service.VideoServiceImpl
	return videoService
}
