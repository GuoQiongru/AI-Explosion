package dao

import (
	"log"
	"TikTok/config"
	"errors"
)

type Like struct {
	Id	int64
	UserId int64
	VideoId int64
	Cancel int8
}

func (Like) TableName() string {
	return "likes"
}

func GetLikeUserIdList(videoId int64) ([]int64, error) {
	var likeUserIdList []int64
	err := Db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.IsLike}).Pluck("user_id", &likeUserIdList).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("get likeUserIdList failed")
	} else {
		return likeUserIdList, nil
	}
}

func UpdateLike(userId int64, videoId int64, actionType int32) error {
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", actionType).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	return nil
}

func InsertLike(likeData Like) error {
	err := Db.Model(Like{}).Create(&likeData).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert data fail")
	}
	return nil
}

func GetLikeInfo(userId int64, videoId int64) (Like, error) {
	var likeInfo Like
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&likeInfo).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return Like{}, nil
		} else {
			log.Println(err.Error())
			return likeInfo, errors.New("get likeInfo failed")
		}
	}
	return likeInfo, nil
}

func GetLikeVideoIdList(userId int64) ([]int64, error) {
	var likeVideoIdList []int64
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "cancel": config.IsLike}).
		Pluck("video_id", &likeVideoIdList).Error
	if err != nil {
		if "record not found" == err.Error() {
			log.Println("there are no likeVideoId")
			return likeVideoIdList, nil
		} else {
			log.Println(err.Error())
			return likeVideoIdList, errors.New("get likeVideoIdList failed")
		}
	}
	return likeVideoIdList, nil
}