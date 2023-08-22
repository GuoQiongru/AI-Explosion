package service

type LikeService interface {
	IsFavourite(videoId int64, userId int64) (bool, error)
	FavouriteCount(videoId int64) (int64, error)
	TotalFavourite(userId int64) (int64, error)
	FavouriteVideoCount(userId int64) (int64, error)

	FavouriteAction(userId int64, videoId int64, actionType int32) error
	GetFavouriteList(userId int64, curId int64) ([]Video, error)
}
