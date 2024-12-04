package stores

import (
	"errors"

	"github.com/lib/pq"
	"github.com/mehakhanaa/complex-micro-blog/models"
	"gorm.io/gorm"
)

type ReplyStore struct {
	db *gorm.DB
}

func (factory *Factory) NewReplyStore() *ReplyStore {
	return &ReplyStore{factory.db}
}

func (store *ReplyStore) CreateReply(uid, commentID uint64, parentReplyID, parentReplyUID *uint64, content string) error {
	newReply := &models.ReplyInfo{
		CommentID:      commentID,
		ParentReplyID:  parentReplyID,
		ParentReplyUID: parentReplyUID,
		Content:        content,
		UID:            uid,
		Like:           pq.Int64Array{},
		Dislike:        pq.Int64Array{},
		IsPublic:       true,
	}

	result := store.db.Create(newReply)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (store *ReplyStore) ValidateReplyExistence(commentID, parentReplyID uint64) (bool, error) {

	var reply models.ReplyInfo
	result := store.db.Where("id = ? AND comment_id = ?", parentReplyID, commentID).First(&reply)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (store *ReplyStore) DeleteReply(uid, replyID uint64) error {
	result := store.db.Model(&models.ReplyInfo{}).Where("id = ? AND uid = ?", replyID, uid).Unscoped().Delete(&models.ReplyInfo{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *ReplyStore) UpdateReply(uid, replyID uint64, content string) error {
	result := store.db.Model(&models.ReplyInfo{}).Where("id = ? AND uid = ?", replyID, uid).Update("content", content)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *ReplyStore) GetReply(replyID uint64) (models.ReplyInfo, error) {
	var reply models.ReplyInfo
	result := store.db.Where("id = ?", replyID).First(&reply)
	if result.Error != nil {
		return models.ReplyInfo{}, result.Error
	}
	return reply, nil
}

func (store *ReplyStore) GetReplyList(commentID uint64) ([]models.ReplyInfo, error) {
	var replyList []models.ReplyInfo
	result := store.db.Where("comment_id = ?", commentID).Order("id desc").Find(&replyList)
	if result.Error != nil {
		return nil, result.Error
	}
	return replyList, nil
}
