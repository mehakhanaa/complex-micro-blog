package stores

import (
	"context"
	"errors"
	"time"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type FollowStore struct {
	db    *gorm.DB
	mongo *mongo.Client
}

func (factory *Factory) NewFollowStore() *FollowStore {
	return &FollowStore{
		factory.db,
		factory.mongo,
	}
}

func (store *FollowStore) FollowUser(uid, followedID uint64) error {

	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "followed_id", Value: followedID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "followed_at", Value: time.Now()},
		}},
	}

	followRecordCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION)
	_, err := followRecordCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))

	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user has liked this followed")
	}
	return err
}

func (store *FollowStore) CancelFollowUser(uid, followedID uint64) error {

	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "followed_id", Value: followedID},
	}

	followRecordCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION)
	_, err := followRecordCollection.DeleteOne(context.Background(), filter)
	if err == mongo.ErrNoDocuments {
		return errors.New("user has not liked this followed")
	}
	return err
}

func (store *FollowStore) GetFollowList(userID uint64) ([]models.FollowInfo, error) {
	var followInfos []models.FollowInfo
	filter := bson.M{
		"uid": userID,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(context.Background(), &followInfos); err != nil {
		return nil, err
	}
	return followInfos, nil
}

func (store *FollowStore) GetFollowedsByUID(uid uint64) (int64, error) {
	filter := bson.M{
		"uid": uid,
	}
	return store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).CountDocuments(context.Background(), filter)
}

func (store *FollowStore) GetFollowerList(userID uint64) ([]models.FollowInfo, error) {
	var followInfos []models.FollowInfo
	filter := bson.M{
		"followed_id": userID,
	}
	cur, err := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(context.Background(), &followInfos); err != nil {
		return nil, err
	}
	return followInfos, nil
}

func (store *FollowStore) GetFollowersByUID(uid uint64) (int64, error) {
	filter := bson.M{
		"followed_id": uid,
	}
	return store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.FOLLOW_RECORD_COLLECTION).CountDocuments(context.Background(), filter)
}
