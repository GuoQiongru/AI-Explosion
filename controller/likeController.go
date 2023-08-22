package controller

import (
	"TikTok/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type likeResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type GetFavouriteListResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg,omitempty"`
	VideoList  []service.Video `json:"video_list,omitempty"`
}

func FavoriteAction(c *gin.Context) {
	strUserId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	strActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(strActionType, 10, 64)
	like := new(service.LikeServiceImpl)
	err := like.FavouriteAction(userId, videoId, int32(actionType))
	if err == nil {
		c.JSON(http.StatusOK, likeResponse{
			StatusCode: 0,
			StatusMsg:  "favourite action success",
		})
	} else {
		log.Printf("like.FavouriteAction Failed: %v", err)
		c.JSON(http.StatusOK, likeResponse{
			StatusCode: 1,
			StatusMsg:  "favourite action fail",
		})
	}
}

func GetFavouriteList(c *gin.Context) {
	strUserId := c.Query("user_id")
	strCurId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	curId, _ := strconv.ParseInt(strCurId, 10, 64)
	like := GetVideo()
	videos, err := like.GetFavouriteList(userId, curId)
	if err == nil {
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 0,
			StatusMsg:  "get favouriteList success",
			VideoList:  videos,
		})
	} else {
		log.Printf("like.GetFavouriteList(userid) Failed: %v", err)
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 1,
			StatusMsg:  "get favouriteList fail ",
		})
	}
}
