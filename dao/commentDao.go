package dao

import (
	"TikTok/config"
	"errors"
	"log"
	"time"
)

type Comment struct {
	Id          int64
	UserId      int64
	VideoId     int64
	CommentText string
	CreateDate  time.Time
	Cancel      int32 //取消评论为1，发布评论为0
}

func (Comment) TableName() string {
	return "comments"
}

func Count(videoId int64) (int64, error) {
	log.Println("CommentDao-Count: running")
	//Init()
	var count int64
	err := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).Count(&count).Error
	if err != nil {
		log.Println("CommentDao-Count: return count failed")
		return -1, errors.New("find comments count failed")
	}
	log.Println("CommentDao-Count: return count success")
	return count, nil
}

func CommentIdList(videoId int64) ([]string, error) {
	var commentIdList []string
	err := Db.Model(Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentIdList).Error
	if err != nil {
		log.Println("CommentIdList:", err)
		return nil, err
	}
	return commentIdList, nil
}

func InsertComment(comment Comment) (Comment, error) {
	log.Println("CommentDao-InsertComment: running")
	err := Db.Model(Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed")
		return Comment{}, errors.New("create comment failed")
	}
	log.Println("CommentDao-InsertComment: return success")
	return comment, nil
}

func DeleteComment(id int64) error {
	log.Println("CommentDao-DeleteComment: running")
	var commentInfo Comment
	result := Db.Model(Comment{}).Where(map[string]interface{}{"id": id, "cancel": config.ValidComment}).First(&commentInfo)
	if result.RowsAffected == 0 {
		log.Println("CommentDao-DeleteComment: return del comment is not exist")
		return errors.New("del comment is not exist")
	}
	err := Db.Model(Comment{}).Where("id = ?", id).Update("cancel", config.InvalidComment).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed")
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success")
	return nil
}

func GetCommentList(videoId int64) ([]Comment, error) {
	log.Println("CommentDao-GetCommentList: running")
	var commentList []Comment
	result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).
		Order("create_date desc").Find(&commentList)
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments")
		return nil, nil
	}
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed")
		return commentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success")
	return commentList, nil
}
