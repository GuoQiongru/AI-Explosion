package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"log"
	"mime/multipart"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type VideoServiceImpl struct {
	UserService
	LikeService
	VideoService
	CommentService
}

func (videoService VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error) {
	videos := make([]Video, 0, config.VideoCount)
	tableVideos, err := dao.GetVideosByLastTime(lastTime)
	if err != nil {
		log.Printf("dao.GetVideosByLastTime(lastTime) failed: %v", err)
		return nil, time.Time{}, err
	}
	err = videoService.copyVideos(&videos, &tableVideos, userId)
	if err != nil {
		log.Printf("videoService.copyVideos(&videos, &tableVideos, userId) failed: %v", err)
		return nil, time.Time{}, err
	}
	nextTime := time.Now()
	if len(tableVideos) > 0 {
		nextTime = tableVideos[len(tableVideos)-1].PublishTime
	}
	return videos, nextTime, nil
}

func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64, title string) error {
	file, err := data.Open()
	if err != nil {
		log.Printf("data.Open() failed: %v", err)
		return err
	}

	videoName := uuid.NewV4().String()
	log.Printf("videoName: %v", videoName)
	err = dao.VideoFTP(file, videoName)
	if err != nil {
		log.Printf("dao.VideoFTP(file, videoName) failed: %v", err)
		return err
	}
	defer file.Close()

	imageName := uuid.NewV4().String()
	ffmpeg.Ffchan <- ffmpeg.Ffmsg{
		videoName,
		imageName,
	}

	err = dao.Save(videoName+".mp4", imageName+".jpg", userId, title)
	if err != nil {
		log.Printf("dao.Save(videoName, imageName, userId) failed: %v", err)
		return err
	}
	return nil
}

func (videoService *VideoServiceImpl) List(userId int64, curId int64) ([]Video, error) {
	data, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		log.Printf("dao.GetVideosByAuthorId(%v) failed: %v", userId, err)
		return nil, err
	}
	result := make([]Video, 0, len(data))
	err = videoService.copyVideos(&result, &data, curId)
	if err != nil {
		log.Printf("videoService.copyVideos(&result, &data, %v)failed: %v", userId, err)
		return nil, err
	}
	return result, nil
}

func (videoService *VideoServiceImpl) copyVideos(result *[]Video, data *[]dao.TableVideo, userId int64) error {
	for _, temp := range *data {
		var video Video
		videoService.creatVideo(&video, &temp, userId)
		*result = append(*result, video)
	}
	return nil
}

func (videoService *VideoServiceImpl) GetVideo(videoId int64, userId int64) (Video, error) {
	var video Video

	data, err := dao.GetVideoByVideoId(videoId)
	if err != nil {
		log.Printf("dao.GetVideoByVideoId(videoId) failed: %v", err)
		return video, err
	}

	videoService.creatVideo(&video, &data, userId)
	return video, nil
}

func (videoService *VideoServiceImpl) creatVideo(video *Video, data *dao.TableVideo, userId int64) {
	var wg sync.WaitGroup
	wg.Add(4)
	var err error
	video.TableVideo = *data
	go func() {
		video.Author, err = videoService.GetUserById(data.AuthorId)
		if err != nil {
			log.Printf("videoService.GetUserByIdWithCurId(data.AuthorId, userId) failed:%v", err)
		}
		wg.Done()
	}()

	go func() {
		video.FavoriteCount, err = videoService.FavouriteCount(data.Id)
		if err != nil {
			log.Printf("videoService.FavouriteCount(data.ID) failed:%v", err)
		}
		wg.Done()
	}()

	go func() {
		video.IsFavorite, err = videoService.IsFavourite(video.Id, userId)
		if err != nil {
			log.Printf("videoService.IsFavourit(video.Id, userId) failed:%v", err)
		}
		wg.Done()
	}()

	go func() {
		video.CommentCount, err = videoService.CountFromVideoId(data.Id)
		if err != nil {
			log.Printf("videoService.CountFromVideoId(data.ID) failed：%v", err)
		} else {
			log.Printf("%d:%d", data.Id, video.CommentCount)
		}
		wg.Done()
	}()

	wg.Wait()
}
