package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ReplyInfo struct {
	gorm.Model
	CommentID      uint64        `gorm:"column:comment_id"`
	ParentReplyID  *uint64       `gorm:"column:reply_to_reply_id"`
	UID            uint64        `gorm:"column:uid"`
	ParentReplyUID *uint64       `gorm:"column:parent_reply_uid"`
	Content        string        `gorm:"column:content"`
	Like           pq.Int64Array `gorm:"column:like;type:bigint[]"`
	Dislike        pq.Int64Array `gorm:"column:dislike;type:bigint[]"`
	IsPublic       bool          `gorm:"column:is_public;default:true"`
}
