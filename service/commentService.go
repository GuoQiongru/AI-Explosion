package service

import (
	"TikTok/dao"
	"time"
)

type CommentService interface {
	CountFromVideoId(id int64) (int64, error)
	Send(comment dao.Comment) (CommentInfo, error)
	DeleteComment(commentId int64) error
	GetList(videoId int64, userId int64) ([]CommentInfo, error)
}

type CommentInfo struct {
	Id         int64  `json:"id"`
	UserInfo   User   `json:"userinfo"`
	Content    string `json:"content"`
	CreateDate string `json:"createDate"`
}

type CommentData struct {
	Id            int64     `json:"id"`
	UserId        int64     `json:"userId"`
	Name          string    `json:"name"`
	FollowCount   int64     `json:"followCount"`
	FollowerCount int64     `json:"followerCount"`
	IsFollow      bool      `json:"isFollow"`
	Content       string    `json:"content"`
	CreateDate    time.Time `json:"createDate"`
}
