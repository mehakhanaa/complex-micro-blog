package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	UserName  string     `gorm:"unique;column:username"`
	NickName  *string    `gorm:"column:nickname"`
	Avatar    string     `gorm:"default:vanilla.webp;column:avatar"`
	Birth     *time.Time `gorm:"column:birth"`
	Gender    *string    `gorm:"column:gender"`
	Authority uint64     `gorm:"default:0;column:authority"`
	Level     uint64     `gorm:"default:1;column:level"`
}

type UserAuthInfo struct {
	gorm.Model
	UID          uint64 `gorm:"unique;column:uid"`
	UserName     string `gorm:"unique;column:username"`
	Salt         string `gorm:"column:salt"`
	PasswordHash string `gorm:"column:psw_hash"`
}

type UserLoginLog struct {
	gorm.Model
	UID         uint64    `gorm:"column:uid"`
	LoginTime   time.Time `gorm:"column:login_time"`
	LoginIP     string    `gorm:"column:login_ip"`
	IsSucceed   bool      `gorm:"column:is_succeed"`
	IfChecked   bool      `gorm:"default:false;column:if_checked"`
	Reason      string    `gorm:"column:reason"`
	Device      string    `gorm:"default:unknown;column:device"`
	Application string    `gorm:"default:unknown;column:application"`
	BearerToken string    `gorm:"column:bearer_token"`
}

type UserPostStatus struct {
	gorm.Model
	UID        uint64        `gorm:"column:uid"`
	Viewed     pq.Int64Array `gorm:"column:viewed;type:bigint[]"`
	Liked      pq.Int64Array `gorm:"column:liked;type:bigint[]"`
	Favourited pq.Int64Array `gorm:"column:favourited;type:bigint[]"`
	Commented  pq.Int64Array `gorm:"column:commented;type:bigint[]"`
}

type UserCommentStatus struct {
	gorm.Model
	UID       uint64        `gorm:"column:uid"`
	Commented pq.Int64Array `gorm:"column:commented;type:bigint[]"`
	Liked     pq.Int64Array `gorm:"column:liked;type:bigint[]"`
	Disliked  pq.Int64Array `gorm:"column:disliked;type:bigint[]"`
}
