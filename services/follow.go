package services

import (
	"github.com/mehakhanaa/complex-micro-blog/models"
	"github.com/mehakhanaa/complex-micro-blog/stores"
)

type FollowService struct {
	followStore *stores.FollowStore
}

func (factory *Factory) NewFollowService() *FollowService {
	return &FollowService{
		followStore: factory.storeFactory.NewFollowStore(),
	}
}

func (service *FollowService) FollowUser(uid, followedID uint64) error {
	return service.followStore.FollowUser(uid, followedID)
}

func (service *FollowService) CancelFollowUser(uid, followedID uint64) error {
	return service.followStore.CancelFollowUser(uid, followedID)
}

func (service *FollowService) GetFollowList(userID uint64) ([]models.FollowInfo, error) {
	return service.followStore.GetFollowList(userID)
}

func (service *FollowService) GetFollowCountByUID(uid uint64) (int64, error) {
	return service.followStore.GetFollowedsByUID(uid)
}

func (service *FollowService) GetFollowerList(userID uint64) ([]models.FollowInfo, error) {
	return service.followStore.GetFollowerList(userID)
}

func (service *FollowService) GetFollowerCountByUID(uid uint64) (int64, error) {
	return service.followStore.GetFollowersByUID(uid)
}
