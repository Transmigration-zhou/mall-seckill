package datamodels

type User struct {
	Id           int64  `json:"id" form:"id" sql:"id"`
	NickName     string `json:"nickName" form:"nickName" sql:"nickName"`
	UserName     string `json:"userName" form:"userName" sql:"userName"`
	HashPassword string `json:"-" form:"password" sql:"password"`
}
