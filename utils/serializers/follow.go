package serializers

import (
	"github.com/mehakhanaa/complex-micro-blog/models"
)

type FollowListResponse struct {
	IDs []uint64 `json:"ids"`
}

func NewFollowListResponse(followInfos []models.FollowInfo) FollowListResponse {
	var ids []uint64
	for _, followInfos := range followInfos {
		ids = append(ids, followInfos.FollowedID)
	}
	return FollowListResponse{IDs: ids}
}

func NewFollowerListResponse(followInfos []models.FollowInfo) FollowListResponse {
	var ids []uint64
	for _, followInfos := range followInfos {
		ids = append(ids, followInfos.UserID)
	}
	return FollowListResponse{IDs: ids}
}
