package stores

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/models"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type PostStore struct {
	db    *gorm.DB
	rds   *redis.Client
	mongo *mongo.Client
}

func (factory *Factory) NewPostStore() *PostStore {
	return &PostStore{
		db:    factory.db,
		rds:   factory.rds,
		mongo: factory.mongo,
	}
}

func (store *PostStore) GetPostList(from string, length int) ([]models.PostInfo, error) {
	var posts []models.PostInfo
	if from != "" {
		if result := store.db.Where("id < ?", from).Order("id desc").Limit(length).Find(&posts); result.Error != nil {
			return nil, result.Error
		}
		return posts, nil
	}
	if result := store.db.Order("id desc").Limit(length).Find(&posts); result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (store *PostStore) GetPostListByUID(uid string) ([]models.PostInfo, error) {
	var userPosts []models.PostInfo
	if result := store.db.Where("uid = ?", uid).Order("id desc").Find(&userPosts); result.Error != nil {
		return nil, result.Error
	}
	return userPosts, nil
}

func (store *PostStore) ValidatePostExistence(postID uint64) (bool, error) {
	var post models.PostInfo
	result := store.db.Where("id = ?", postID).First(&post)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (store *PostStore) GetPostInfo(postID uint64) (models.PostInfo, int64, int64, error) {
	post := models.PostInfo{}
	result := store.db.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return models.PostInfo{}, 0, 0, result.Error
	}

	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	likeCount, err := postLikeCollection.CountDocuments(context.Background(), bson.D{{Key: "post_id", Value: postID}})
	if err != nil {
		return models.PostInfo{}, 0, 0, err
	}
	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	favouriteCount, err := postFavouriteCollection.CountDocuments(context.Background(), bson.D{{Key: "post_id", Value: postID}})
	if err != nil {
		return models.PostInfo{}, 0, 0, err
	}

	return post, likeCount, favouriteCount, nil
}

func (store *PostStore) CreatePost(uid uint64, ipAddr string, postReqData types.PostCreateBody) (models.PostInfo, error) {
	var imageFileNames []string

	for _, imageUUID := range postReqData.Images {
		srcImage, err := os.Open(filepath.Join(consts.POST_IMAGE_CACHE_PATH, imageUUID+".webp"))
		if err != nil {
			return models.PostInfo{}, err
		}
		defer srcImage.Close()
		dstImage, err := os.Create(filepath.Join(consts.POST_IMAGE_PATH, imageUUID+".webp"))
		if err != nil {
			return models.PostInfo{}, err
		}
		defer dstImage.Close()
		_, err = io.Copy(dstImage, srcImage)
		if err != nil {
			return models.PostInfo{}, err
		}
		imageFileNames = append(imageFileNames, imageUUID+".webp")

		ctx := context.Background()
		tx := store.rds.TxPipeline()

		_, err = tx.XAdd(ctx, &redis.XAddArgs{
			Stream: consts.CACHE_IMG_CLEAN_STREAM,
			Values: map[string]interface{}{"filename": imageUUID + ".webp"},
		}).Result()
		if err != nil {
			tx.Discard()
			return models.PostInfo{}, err
		}

		var sb strings.Builder
		sb.WriteString(consts.CACHE_IMAGE_LIST)
		sb.WriteString(":")
		sb.WriteString(imageUUID)
		fmt.Println(sb.String())
		_, err = tx.Del(ctx, sb.String()).Result()
		if err != nil {
			tx.Discard()
			return models.PostInfo{}, err
		}

		_, err = tx.Exec(ctx)
		if err != nil {
			tx.Discard()
			return models.PostInfo{}, err
		}
	}

	postInfo := models.PostInfo{
		ParentPostID: nil,
		UID:          uid,
		IpAddrress:   &ipAddr,
		Title:        postReqData.Title,
		Content:      postReqData.Content,
		Images:       imageFileNames,
		Like:         pq.Int64Array{},
		Favourite:    pq.Int64Array{},
		Farward:      pq.Int64Array{},
		IsPublic:     true,
	}
	result := store.db.Create(&postInfo)
	return postInfo, result.Error
}

func (store *PostStore) CachePostImage(image []byte) (string, error) {

	var (
		fileNameBuilder strings.Builder
		UUID            string
		savePath        string
	)
	for {
		UUID = strings.ReplaceAll(uuid.New().String(), "-", "")
		fileNameBuilder.WriteString(UUID)
		fileNameBuilder.WriteString(".webp")
		savePath = filepath.Join(consts.POST_IMAGE_CACHE_PATH, fileNameBuilder.String())
		_, err := os.Stat(savePath)
		if os.IsExist(err) {
			fileNameBuilder.Reset()
			continue
		}
		break
	}

	file, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewReader(image))
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	var sb strings.Builder
	sb.WriteString(consts.CACHE_IMAGE_LIST)
	sb.WriteString(":")
	sb.WriteString(UUID)

	_, err = store.rds.HSet(ctx, sb.String(), map[string]interface{}{
		"filename": fileNameBuilder.String(),
		"expire":   time.Now().Add(consts.CACHE_IMAGE_EXPIRE_TIME * time.Second).Unix(),
	}).Result()

	if err != nil {
		return "", err
	}

	return UUID, nil
}

func (store *PostStore) CheckCacheImageAvaliable(uuid string) (bool, error) {

	ctx := context.Background()

	var sb strings.Builder
	sb.WriteString(consts.CACHE_IMAGE_LIST)
	sb.WriteString(":")
	sb.WriteString(uuid)

	flag := false
	keys, err := store.rds.Keys(ctx, consts.CACHE_IMAGE_LIST+":*").Result()
	if err != nil {
		return false, err
	}
	for _, key := range keys {
		if store.rds.HGet(ctx, key, "filename").Val() == uuid+".webp" {

			expire, err := store.rds.HGet(ctx, key, "expire").Int64()
			if err != nil {
				return false, err
			}

			if time.Now().Unix() > expire {
				return false, nil
			}
			flag = true
			break
		}
	}

	if !flag {
		return false, nil
	}

	_, err = os.Stat(filepath.Join(consts.POST_IMAGE_CACHE_PATH, uuid+".webp"))

	if os.IsNotExist(err) {

		tx := store.rds.TxPipeline()
		var sb strings.Builder
		sb.WriteString(consts.CACHE_IMAGE_LIST)
		sb.WriteString(":")
		sb.WriteString(uuid)
		_, err = tx.Del(ctx, sb.String()).Result()
		if err != nil {
			tx.Discard()
			return false, err
		}
		_, err = tx.Exec(ctx)
		if err != nil {
			tx.Discard()
			return false, err
		}
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (store *PostStore) LikePost(uid, postID int64) error {

	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "liked_at", Value: time.Now()},
		}},
	}

	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	_, err := postLikeCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))

	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user has liked this post")
	}
	return err
}

func (store *PostStore) CancelLikePost(uid, postID int64) error {
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}

	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	_, err := postLikeCollection.DeleteOne(context.Background(), filter)
	if mongo.ErrNoDocuments == err {
		return errors.New("user has not liked this post")
	}
	return err
}

func (store *PostStore) FavouritePost(uid, postID int64) error {
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "favourited_at", Value: time.Now()},
		}},
	}

	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	_, err := postFavouriteCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("user has favourited this post")
	}
	return err
}

func (store *PostStore) CancelFavouritePost(uid, postID int64) error {
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}

	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	_, err := postFavouriteCollection.DeleteOne(context.Background(), filter)
	if mongo.ErrNoDocuments == err {
		return errors.New("user has not favourited this post")
	}
	return err
}

func (store *PostStore) GetPostUserStatus(uid, postID int64) (bool, bool, error) {

	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "post_id", Value: postID},
	}

	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	count, err := postLikeCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, false, err
	}
	isLiked := count > 0

	postFavouriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	count, err = postFavouriteCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, false, err
	}
	isFavourited := count > 0

	return isLiked, isFavourited, nil
}

func (store *PostStore) DeletePost(postID uint64) error {
	return store.db.Where("id = ?", postID).Unscoped().Delete(&models.PostInfo{}).Error
}
