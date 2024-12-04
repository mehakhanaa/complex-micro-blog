package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CommentInfo struct {
	gorm.Model
	PostID   uint64        `gorm:"column:post_id"`
	UID      uint64        `gorm:"column:uid"`
	Username string        `gorm:"column:username"`
	Content  string        `gorm:"column:content"`
	Like     pq.Int64Array `gorm:"column:like;type:bigint[]"`
	Dislike  pq.Int64Array `gorm:"column:dislike;type:bigint[]"`
	IsPublic bool          `gorm:"column:is_public;default:true"`
}
