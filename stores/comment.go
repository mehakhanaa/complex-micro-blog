package stores

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type CommentStore struct {
	db    *gorm.DB
	mongo *mongo.Client
}

func (factory *Factory) NewCommentStore() *CommentStore {
	return &CommentStore{
		factory.db,
		factory.mongo,
	}
}

func (store *CommentStore) CreateComment(uid uint64, username string, postID uint64, content string) (uint64, error) {
	newComment := models.CommentInfo{
		PostID:   postID,
		Username: username,
		Content:  content,
		UID:      uid,
		Like:     pq.Int64Array{},
		Dislike:  pq.Int64Array{},
		IsPublic: true,
	}

	result := store.db.Create(&newComment)
	if result.Error != nil {
		return 0, result.Error
	}
	return uint64(newComment.ID), nil
}

func (store *CommentStore) ValidateCommentExistence(commentID uint64) (bool, error) {
	var comment models.CommentInfo
	result := store.db.Where("id = ?", commentID).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (store *CommentStore) UpdateComment(commentID uint64, content string) error {
	commentInfo := new(models.CommentInfo)
	result := store.db.Where("id = ?", commentID).First(commentInfo)
	if result.Error != nil {
		return result.Error
	}

	commentInfo.Content = content
	result = store.db.Save(commentInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *CommentStore) DeleteComment(commentID uint64) error {
	return store.db.Where("id = ?", commentID).Unscoped().Delete(&models.CommentInfo{}).Error
}

func (store *CommentStore) GetCommentList(postID uint64) ([]models.CommentInfo, error) {
	var commentInfos []models.CommentInfo
	result := store.db.Where("post_id = ?", postID).Order("id desc").Find(&commentInfos)
	if result.Error != nil {
		return nil, result.Error
	}
	return commentInfos, nil
}

func (store *CommentStore) GetCommentInfo(commentID uint64) (models.CommentInfo, int64, error) {
	comment := models.CommentInfo{}
	result := store.db.Where("id = ?", commentID).First(&comment)
	if result.Error != nil {
		return comment, 0, result.Error
	}

	commentRateCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.COMMENT_RATE_COLLECTION)
	filter := bson.D{
		{Key: "comment_id", Value: commentID},
		{Key: "rate", Value: "like"},
	}
	ctx := context.Background()
	defer ctx.Done()

	likeCount, err := commentRateCollection.CountDocuments(ctx, filter)
	if err != nil {
		return comment, 0, err
	}

	return comment, likeCount, nil
}

func (store *CommentStore) GetCommentUserStatus(uid, commentID uint64) (bool, bool, error) {
	commentRateCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.COMMENT_RATE_COLLECTION)
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "comment_id", Value: commentID},
		{Key: "rate", Value: "like"},
	}
	ctx := context.Background()
	defer ctx.Done()

	count, err := commentRateCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, false, err
	}
	isLiked := count > 0

	filter[2].Value = "dislike"
	count, err = commentRateCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, false, err
	}
	isDisliked := count > 0

	return isLiked, isDisliked, nil
}

func (store *CommentStore) LikeComment(uid, commentID uint64) error {
	commentRateCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.COMMENT_RATE_COLLECTION)
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "comment_id", Value: commentID},
	}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "rate", Value: "like"},
				{Key: "rated_at", Value: time.Now()},
			},
		},
	}

	_, err := commentRateCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))

	return err
}

func (store *CommentStore) CancelLikeComment(uid, commentID uint64) error {
	commentRateCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.COMMENT_RATE_COLLECTION)
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "comment_id", Value: commentID},
		{Key: "rate", Value: "like"},
	}

	result, err := commentRateCollection.DeleteOne(context.Background(), filter)
	if result.DeletedCount == 0 {
		return errors.New("user has not liked this comment")
	}
	return err
}

func (store *CommentStore) DislikeComment(uid, commentID uint64) error {
	commentRateCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.COMMENT_RATE_COLLECTION)
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "comment_id", Value: commentID},
	}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "rate", Value: "dislike"},
				{Key: "rated_at", Value: time.Now()},
			},
		},
	}

	_, err := commentRateCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))

	return err
}

func (store *CommentStore) CancelDislikeComment(uid, commentID uint64) error {
	commentRateCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.COMMENT_RATE_COLLECTION)
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "comment_id", Value: commentID},
		{Key: "rate", Value: "dislike"},
	}

	result, err := commentRateCollection.DeleteOne(context.Background(), filter)
	if result.DeletedCount == 0 {
		return errors.New("user has not disliked this comment")
	}

	return err
}
