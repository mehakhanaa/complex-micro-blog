package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PostInfo struct {
	gorm.Model
	ParentPostID *uint64        `gorm:"column:parent_post_id"`
	UID          uint64         `gorm:"column:uid"`
	IpAddrress   *string        `gorm:"column:ip_address"`
	Title        string         `gorm:"column:title"`
	Content      string         `gorm:"column:content"`
	Images       pq.StringArray `gorm:"column:images;type:text[]"`
	Like         pq.Int64Array  `gorm:"column:like;type:bigint[]"`
	Favourite    pq.Int64Array  `gorm:"column:favourite;type:bigint[]"`
	Farward      pq.Int64Array  `gorm:"column:farward;type:bigint[]"`
	IsPublic     bool           `gorm:"column:is_public;default:true"`
}
