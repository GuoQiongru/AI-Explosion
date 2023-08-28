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
	UserInfo   User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type CommentData struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	Name       string    `json:"name"`
	Avatar     string    `json:"avatar"`
	Content    string    `json:"content"`
	CreateDate time.Time `json:"create_date"`
}
