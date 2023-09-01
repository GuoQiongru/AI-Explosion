package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware/redis"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LikeServiceImpl struct {
	VideoService
	UserService
	LikeService
}

func (like *LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			log.Printf("IsFavourite RedisLikeUserId query key failed: %v", err)
			return false, err
		}
		exist, err1 := redis.RdbLikeUserId.SIsMember(redis.Ctx, strUserId, videoId).Result()
		if err1 != nil {
			log.Printf("IsFavourite RedisLikeUserId query value failed: %v", err1)
			return false, err1
		}
		return exist, nil
	} else {
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			if err != nil {
				log.Printf("IsFavourite RedisLikeVideoId query key failed: %v", err)
				return false, err
			}
			exist, err1 := redis.RdbLikeVideoId.SIsMember(redis.Ctx, strVideoId, userId).Result()
			if err1 != nil {
				log.Printf("IsFavourite RedisLikeVideoId query value failed: %v", err1)
				return false, err1
			}
			return exist, nil
		} else {
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("IsFavourite RedisLikeUserId add value failed")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return false, err
			}
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("IsFavourite RedisLikeUserId add time failed")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return false, err
			}
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			if err1 != nil {
				log.Print(err1.Error())
				return false, err1
			}
			for _, likeVideoId := range videoIdList {
				redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId)
			}
			exist, err2 := redis.RdbLikeUserId.SIsMember(redis.Ctx, strUserId, videoId).Result()
			if err2 != nil {
				log.Printf("IsFavourite RedisLikeUserId query value failed: %v", err2)
				return false, err2
			}
			return exist, nil
		}
	}
}

func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	strVideoId := strconv.FormatInt(videoId, 10)
	if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
		if err != nil {
			log.Printf("FavouriteCount RedisLikeVideoId query key failed: %v", err)
			return 0, err
		}
		count, err1 := redis.RdbLikeVideoId.SCard(redis.Ctx, strVideoId).Result()
		if err1 != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId query count failed: %v", err1)
			return 0, err1
		}
		return count - 1, nil
	} else {
		if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
			log.Printf("FavouriteCount RedisLikeVideoId add value failed: ")
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId,
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			log.Printf("FavouriteCount RedisLikeVideoId add time failed")
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Print(err1.Error())
			return 0, err1
		}
		for _, likeUserId := range userIdList {
			redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId)
		}
		count, err2 := redis.RdbLikeVideoId.SCard(redis.Ctx, strVideoId).Result()
		if err2 != nil {
			log.Printf("FavouriteCount RedisLikeVideoId query count failed: %v", err2)
			return 0, err2
		}
		return count - 1, nil
	}
}

func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, actionType int32) error {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	sb := strings.Builder{}
	sb.WriteString(strUserId)
	sb.WriteString(" ")
	sb.WriteString(strVideoId)

	if actionType == config.LikeAction {
		if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
			if err != nil {
				log.Printf("FavouriteAction RedisLikeUserId query key failed: %v", err)
				return err
			}
			if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
				log.Printf("FavouriteAction RedisLikeUserId add value failed: %v", err1)
				return err1
			} else {
				like.addLike(userId, videoId)
			}
		} else {
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("FavouriteAction RedisLikeUserId add value failed: ")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("FavouriteAction RedisLikeUserId add time failed: ")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			if err1 != nil {
				return err1
			}
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					log.Printf("FavouriteAction RedisLikeUserId add value failed: ")
					redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				log.Printf("FavouriteAction RedisLikeUserId add value failed: %v", err2)
				return err2
			} else {
				like.addLike(userId, videoId)
			}
		}
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			if err != nil {
				log.Printf("FavouriteAction RedisLikeVideoId query key failed: %v", err)
				return err
			}
			if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
				log.Printf("FavouriteAction RedisLikeVideoId add value failed: %v", err1)
				return err1
			}
		} else {
			if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("FavouriteAction RedisLikeVideoId add value failed")
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("FavouriteAction RedisLikeVideoId add time failed: ")
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			userIdList, err1 := dao.GetLikeUserIdList(videoId)
			if err1 != nil {
				return err1
			}
			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					log.Printf("FavouriteAction RedisLikeVideoId add value failed: ")
					redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			if _, err2 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				log.Printf("FavouriteAction RedisLikeVideoId add value failed : %v", err2)
				return err2
			}
		}
	} else {
		if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
			if err != nil {
				log.Printf("FavouriteAction RedisLikeUserId query key failed: %v", err)
				return err
			}
			if _, err1 := redis.RdbLikeUserId.SRem(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
				log.Printf("FavouriteAction RedisLikeUserId del value failed: %v", err1)
				return err1
			} else {
				like.delLike(userId, videoId)
			}
		} else {
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("FavouriteAction RedisLikeUserId add value failed: ")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("FavouriteAction RedisLikeUserId add time failed: ")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			if err1 != nil {
				return err1
			}
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					log.Printf("FavouriteAction RedisLikeUserId add value failed: ")
					redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			if _, err2 := redis.RdbLikeUserId.SRem(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				log.Printf("FavouriteAction RedisLikeUserId del value failed: %v", err2)
				return err2
			} else {
				like.delLike(userId, videoId)
			}
		}

		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			if err != nil {
				return err
			}
			if _, err1 := redis.RdbLikeVideoId.SRem(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
				return err1
			}
		} else {
			if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}

			userIdList, err1 := dao.GetLikeUserIdList(videoId)
			if err1 != nil {
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err1
			}

			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			if _, err2 := redis.RdbLikeVideoId.SRem(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func (like *LikeServiceImpl) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	strUserId := strconv.FormatInt(userId, 10)
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			return nil, err
		}
		videoIdList, err1 := redis.RdbLikeUserId.SMembers(redis.Ctx, strUserId).Result()
		if err1 != nil {
			return nil, err1
		}
		favoriteVideoList := new([]Video)
		i := len(videoIdList) - 1
		if i == 0 {
			return *favoriteVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			videoId, _ := strconv.ParseInt(videoIdList[j], 10, 64)
			if videoId == config.DefaultRedisValue {
				continue
			}
			go like.addFavouriteVideoList(videoId, curId, favoriteVideoList, &wg)
		}
		wg.Wait()
		return *favoriteVideoList, nil
	} else {
		if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err
		}
		_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err
		}
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println(err1.Error())
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err1
		}

		for _, likeVideoId := range videoIdList {
			if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err2 != nil {
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return nil, err2
			}
		}
		favoriteVideoList := new([]Video)
		i := len(videoIdList) - 1
		if i == 0 {
			return *favoriteVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			if videoIdList[j] == config.DefaultRedisValue {
				continue
			}
			go like.addFavouriteVideoList(videoIdList[j], curId, favoriteVideoList, &wg)
		}
		wg.Wait()
		return *favoriteVideoList, nil
	}
}

func (like *LikeServiceImpl) addFavouriteVideoList(videoId int64, curId int64, favoriteVideoList *[]Video, wg *sync.WaitGroup) {
	defer wg.Done()
	video, err := like.GetVideo(videoId, curId)
	if err != nil {
		log.Println(errors.New("this favourite video is miss"))
		return
	}
	*favoriteVideoList = append(*favoriteVideoList, video)
}

func (like *LikeServiceImpl) TotalFavourite(userId int64) (int64, error) {
	videoIdList, err := like.GetVideoIdList(userId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}
	var sum int64
	videoLikeCountList := new([]int64)
	i := len(videoIdList)
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j < i; j++ {
		go like.addVideoLikeCount(videoIdList[j], videoLikeCountList, &wg)
	}
	wg.Wait()
	for _, count := range *videoLikeCountList {
		sum += count
	}
	return sum, nil
}

func (like *LikeServiceImpl) FavouriteVideoCount(userId int64) (int64, error) {
	strUserId := strconv.FormatInt(userId, 10)
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		if err != nil {
			return 0, err
		} else {
			count, err1 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
			if err1 != nil {
				return 0, err1
			}
			return count - 1, nil

		}
	} else {
		if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return 0, err
		}
		_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return 0, err
		}
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Print(err1.Error())
			return 0, err1
		}
		for _, likeVideoId := range videoIdList {
			if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return 0, err1
			}
		}
		count, err2 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
		if err2 != nil {
			return 0, err2
		}
		return count - 1, nil
	}
}

func (like *LikeServiceImpl) addVideoLikeCount(videoId int64, videoLikeCountList *[]int64, wg *sync.WaitGroup) {
	defer wg.Done()
	count, err := like.FavouriteCount(videoId)
	if err != nil {
		log.Print(err.Error())
		return
	}
	*videoLikeCountList = append(*videoLikeCountList, count)
}

func GetLikeService() LikeServiceImpl {
	var userService UserServiceImpl
	var videoService VideoServiceImpl
	var likeService LikeServiceImpl
	userService.LikeService = &likeService
	likeService.VideoService = &videoService
	videoService.UserService = &userService
	return likeService
}

func (like *LikeServiceImpl) delLike(userId int64, videoId int64) bool {
	flag := false
	likeInfo, err := dao.GetLikeInfo(userId, videoId)
	if err != nil {
		log.Print(err.Error())
		flag = true
	} else {
		if likeInfo == (dao.Like{}) {
			log.Print(errors.New("can't find data,this action invalid").Error())
		} else {
			if err := dao.UpdateLike(userId, videoId, config.Unlike); err != nil {
				log.Print(err.Error())
				flag = true
			}
		}
	}
	if !flag {
		log.Println("删除失败")
		return false
	}
	return true
}

func (like *LikeServiceImpl) addLike(userId int64, videoId int64) bool {
	flag := false
	var likeData dao.Like
	likeInfo, err := dao.GetLikeInfo(userId, videoId)
	if err != nil {
		log.Print(err.Error())
		flag = true
	} else {
		if likeInfo == (dao.Like{}) {
			likeData.UserId = userId
			likeData.VideoId = videoId
			likeData.Cancel = config.IsLike
			if err := dao.InsertLike(likeData); err != nil {
				log.Print(err.Error())
				flag = true
			}
		} else {
			if err := dao.UpdateLike(userId, videoId, config.IsLike); err != nil {
				log.Print(err.Error())
				flag = true
			}
		}
	}
	if !flag {
		log.Println("添加失败")
		return false
	}
	return true
}
