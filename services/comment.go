package services

import (
	"errors"

	"github.com/mehakhanaa/complex-micro-blog/models"
	"github.com/mehakhanaa/complex-micro-blog/stores"
)

type CommentService struct {
	commentStore *stores.CommentStore
}

func (factory *Factory) NewCommentService() *CommentService {
	return &CommentService{
		commentStore: factory.storeFactory.NewCommentStore(),
	}
}

func (service *CommentService) CreateComment(uid uint64, postID uint64, content string, postStore *stores.PostStore, userStore *stores.UserStore) (uint64, error) {

	existance, err := postStore.ValidatePostExistence(postID)
	if err != nil {
		return 0, err
	}
	if !existance {
		return 0, errors.New("post does not exist")
	}

	user, err := userStore.GetUserByUID(uid)
	if err != nil {
		return 0, err
	}

	commentID, err := service.commentStore.CreateComment(uid, user.UserName, postID, content)
	if err != nil {
		return 0, err
	}
	return commentID, nil
}

func (service *CommentService) UpdateComment(commentID uint64, content string) error {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	err = service.commentStore.UpdateComment(commentID, content)
	if err != nil {
		return err
	}

	return nil
}

func (service *CommentService) DeleteComment(commentID uint64) error {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	err = service.commentStore.DeleteComment(commentID)
	if err != nil {

		return err
	}

	return nil
}

func (service *CommentService) GetCommentList(postID uint64) ([]models.CommentInfo, error) {
	return service.commentStore.GetCommentList(postID)
}

func (service *CommentService) GetCommentInfo(commentID uint64) (models.CommentInfo, int64, error) {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return models.CommentInfo{}, 0, err
	}
	if !exists {
		return models.CommentInfo{}, 0, errors.New("comment does not exist")
	}

	return service.commentStore.GetCommentInfo(commentID)
}

func (service *CommentService) GetCommentUserStatus(uid, commentID uint64) (bool, bool, error) {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return false, false, err
	}
	if !exists {
		return false, false, errors.New("comment does not exist")
	}

	isLiked, isDisliked, err := service.commentStore.GetCommentUserStatus(uid, commentID)
	if err != nil {
		return false, false, err
	}

	return isLiked, isDisliked, nil
}

func (service *CommentService) LikeComment(uid, commentID uint64) error {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	err = service.commentStore.LikeComment(uid, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (service *CommentService) CancelLikeComment(uid, commentID uint64) error {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	err = service.commentStore.CancelLikeComment(uid, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (service *CommentService) DislikeComment(uid, commentID uint64) error {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	err = service.commentStore.DislikeComment(uid, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (service *CommentService) CancelDislikeComment(uid, commentID uint64) error {

	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	err = service.commentStore.CancelDislikeComment(uid, commentID)
	if err != nil {
		return err
	}

	return nil
}
