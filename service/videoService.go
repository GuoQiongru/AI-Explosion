package service

import (
	"TikTok/dao"
	"mime/multipart"
	"time"
)

type Video struct {
	dao.TableVideo
	Author        User  `json:"author"`
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount  int64 `json:"comment_count"`
	IsFavorite    bool  `json:"is_favorite"`
}

type VideoService interface {
	Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error)

	GetVideo(videoId int64, userId int64) (Video, error)

	Publish(data *multipart.FileHeader, userId int64, title string) error

	List(userId int64, curId int64) ([]Video, error)

	GetVideoIdList(userId int64) ([]int64, error)
}
