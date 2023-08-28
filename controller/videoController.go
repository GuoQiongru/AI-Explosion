package controller

import (
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
	VideoList []service.Video `json:"video_list"`
	NextTime  int64           `json:"next_time"`
}

func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	var lastTime time.Time
	if inputTime != "" && inputTime != "0" {
		me, _ := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.Unix(me, 0)
	} else {
		lastTime = time.Now()
	}
	log.Printf("LastTime: %v", lastTime)
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)

	videoService := GetVideo()
	feed, nextTime, err := videoService.Feed(lastTime, userId)
	if err != nil {
		log.Printf("videoService.Feed(lastTime, userId) Failed: %v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "feed Failed"},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

func Publish(c *gin.Context) {
	file, err := c.FormFile("data")
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	fmt.Printf("Video Publish UserId:%v\n", userId)
	title := c.PostForm("title")
	log.Printf("Video Publish title:%v\n", title)
	if err != nil {
		log.Printf("Video Publish Failed :%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoService := GetVideo()
	err = videoService.Publish(file, userId, title)
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
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)

	videoService := GetVideo()
	feed, err := videoService.List(userId, curId)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "publishList failed"},
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
	var userService service.UserServiceImpl
	var videoService service.VideoServiceImpl
	var likeService service.LikeServiceImpl
	var commentService service.CommentServiceImpl
	userService.LikeService = &likeService
	likeService.VideoService = &videoService
	videoService.CommentService = &commentService
	videoService.LikeService = &likeService
	videoService.UserService = &userService
	return videoService
}
