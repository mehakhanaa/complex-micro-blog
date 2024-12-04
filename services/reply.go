package services

import (
	"errors"

	"github.com/mehakhanaa/complex-micro-blog/models"
	"github.com/mehakhanaa/complex-micro-blog/stores"
)

type ReplyService struct {
	replyStore *stores.ReplyStore
}

func (factory *Factory) NewReplyService() *ReplyService {
	return &ReplyService{
		replyStore: factory.storeFactory.NewReplyStore(),
	}
}

func (service *ReplyService) CreateReply(uid, commentID, parentReplyID uint64, content string, commentStore *stores.CommentStore, userStore *stores.UserStore) error {

	isExist, err := commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("comment does not exist")
	}

	var parentReplyUIDField *uint64 = nil

	if parentReplyID != 0 {
		isExist, err := service.replyStore.ValidateReplyExistence(commentID, parentReplyID)
		if err != nil {
			return err
		}
		if !isExist {
			return errors.New("reply does not exist")
		}
		parentReplyInfo, err := service.replyStore.GetReply(parentReplyID)
		if err != nil {
			return err
		}
		parentReplyUIDField = &parentReplyInfo.UID
	}

	var parentReplyIDField *uint64 = nil
	if parentReplyID != 0 {
		parentReplyIDField = &parentReplyID
	}

	err = service.replyStore.CreateReply(uid, commentID, parentReplyIDField, parentReplyUIDField, content)
	if err != nil {
		return err
	}
	return nil
}

func (service *ReplyService) DeleteReply(uid, replyID uint64) error {

	err := service.replyStore.DeleteReply(uid, replyID)
	if err != nil {

		return err
	}

	return nil
}

func (service *ReplyService) UpdateReply(uid, replyID uint64, content string) error {

	err := service.replyStore.UpdateReply(uid, replyID, content)
	if err != nil {
		return err
	}

	return nil
}

func (service *ReplyService) GetReplyList(commentID uint64) ([]uint64, error) {

	replyList, err := service.replyStore.GetReplyList(commentID)
	if err != nil {
		return nil, err
	}
	replyListUint64 := make([]uint64, len(replyList))
	for index, reply := range replyList {
		replyListUint64[index] = uint64(reply.ID)
	}

	return replyListUint64, nil
}

func (service *ReplyService) GetReplyDetail(replyID uint64) (models.ReplyInfo, error) {

	reply, err := service.replyStore.GetReply(replyID)
	if err != nil {
		return models.ReplyInfo{}, err
	}

	return reply, nil
}
