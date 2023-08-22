package service

import "TikTok/dao"

type UserService interface {
	GetTableUserList() []dao.TableUser

	GetTableUserByUsername(name string) dao.TableUser

	GetTableUserById(id int64) dao.TableUser

	InsertTableUser(tableUser *dao.TableUser) bool

	GetUserById(id int64) (User, error)
}

type User struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count,omitempty"`
	FollowerCount  int64  `json:"follower_count,omitempty"`
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
	Avatar         string `json:"avatar"`
}
