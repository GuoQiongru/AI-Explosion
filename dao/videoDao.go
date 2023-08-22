package dao

import (
	"TikTok/config"
	"TikTok/middleware/ftp"
	"io"
	"log"
	"time"
)

type TableVideo struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"`
}

func (TableVideo) TableName() string {
	return "videos"
}

func GetVideosByAuthorId(authorId int64) ([]TableVideo, error) {
	var data []TableVideo
	result := Db.Where(&TableVideo{AuthorId: authorId}).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func GetVideoByVideoId(videoId int64) (TableVideo, error) {
	var tableVideo TableVideo
	tableVideo.Id = videoId
	result := Db.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil

}

func GetVideosByLastTime(lastTime time.Time) ([]TableVideo, error) {
	var videos []TableVideo
	result := Db.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

func Save(videoName string, imageName string, authorId int64, title string) error {
	var video TableVideo
	video.PublishTime = time.Now()
	video.PlayUrl = config.UrlPrefix + videoName
	video.CoverUrl = config.UrlPrefix + imageName
	video.AuthorId = authorId
	video.Title = title
	result := Db.Save(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetVideoIdsByAuthorId(authorId int64) ([]int64, error) {
	var id []int64
	result := Db.Model(&TableVideo{}).Where("author_id", authorId).Pluck("id", &id)
	if result.Error != nil {
		return nil, result.Error
	}
	return id, nil
}

func VideoFTP(file io.Reader, videoName string) error {
	ftp.MyFTP.Cwd("video")
	err := ftp.MyFTP.Stor(videoName+".mp4", file)
	if err != nil {
		log.Println("VideoFTP Failed")
		return err
	}
	log.Println("VideoFTP Successfully")
	return nil
}

func ImageFTP(file io.Reader, imageName string) error {
	ftp.MyFTP.Cwd("images")
	if err := ftp.MyFTP.Stor(imageName, file); err != nil {
		log.Println("ImageFTP Failed")
		return err
	}
	log.Println("ImageFTP Successfully")
	return nil
}
