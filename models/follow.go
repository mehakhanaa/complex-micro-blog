package models

import (
	"time"
)

type FollowInfo struct {
	UserID     uint64    `bson:"uid"`
	FollowedID uint64    `bson:"followed_id"`
	FollowedAt time.Time `bson:"followed_at"`
}
