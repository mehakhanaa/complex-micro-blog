package stores

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"

	"github.com/lib/pq"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/models"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	db    *gorm.DB
	rds   *redis.Client
	mongo *mongo.Client
}

func (factory *Factory) NewUserStore() *UserStore {
	return &UserStore{
		factory.db,
		factory.rds,
		factory.mongo,
	}
}

func (store *UserStore) RegisterUserByUsername(username string, salt string, hashedPassword string) error {
	tx := store.db.Begin()

	user := models.UserInfo{
		UserName: username,
		NickName: &username,
	}
	result := tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	uid := user.ID
	userAuthInfo := models.UserAuthInfo{
		UID:          uint64(uid),
		UserName:     username,
		Salt:         salt,
		PasswordHash: hashedPassword,
	}
	result = tx.Create(&userAuthInfo)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	userPostStatus := models.UserPostStatus{
		UID:        uint64(uid),
		Viewed:     pq.Int64Array{},
		Liked:      pq.Int64Array{},
		Favourited: pq.Int64Array{},
		Commented:  pq.Int64Array{},
	}
	result = tx.Create(&userPostStatus)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	userCommentStatus := models.UserCommentStatus{
		UID:       uint64(uid),
		Liked:     pq.Int64Array{},
		Disliked:  pq.Int64Array{},
		Commented: pq.Int64Array{},
	}
	result = tx.Create(&userCommentStatus)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func (store *UserStore) GetUserByUID(uid uint64) (*models.UserInfo, error) {
	user := new(models.UserInfo)
	result := store.db.Where("id = ?", uid).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (store *UserStore) GetUserByUsername(username string) (*models.UserInfo, error) {
	user := new(models.UserInfo)
	result := store.db.Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (store *UserStore) GetUserAuthInfoByUsername(username string) (*models.UserAuthInfo, error) {
	userAuthInfo := new(models.UserAuthInfo)
	result := store.db.Where("username = ?", username).First(userAuthInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return userAuthInfo, nil
}

func (store *UserStore) CreateUserLoginLog(userLoginLogInfo *models.UserLoginLog) error {
	result := store.db.Create(userLoginLogInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (store *UserStore) CreateUserAvaliableToken(token string, claims *types.BearerTokenClaims) error {
	var sb strings.Builder
	sb.WriteString(consts.REDIS_AVAILABLE_USER_TOKEN_LIST)
	sb.WriteRune(':')
	sb.WriteString(strconv.FormatUint(claims.UID, 10))
	key := sb.String()

	ctx := context.Background()

	length, err := store.rds.LLen(ctx, key).Result()
	if err != nil {
		return err
	}

	fmt.Println(length)

	tx := store.rds.TxPipeline()

	if length >= consts.MAX_TOKENS_PER_USER {

		_, err = tx.LTrim(ctx, key, length-4, -1).Result()
		if err != nil {
			tx.Discard()
			return err
		}
	}

	_, err = tx.RPush(ctx, key, token).Result()
	if err != nil {
		tx.Discard()
		return err
	}

	_, err = tx.Exec(ctx)
	if err != nil {
		tx.Discard()
		return err
	}

	return nil
}

func (store *UserStore) BanUserToken(uid uint64, token string) error {
	var sb strings.Builder
	sb.WriteString(consts.REDIS_AVAILABLE_USER_TOKEN_LIST)
	sb.WriteRune(':')
	sb.WriteString(strconv.FormatUint(uid, 10))
	key := sb.String()

	ctx := context.Background()
	tx := store.rds.TxPipeline()

	_, err := tx.LRem(ctx, key, 0, token).Result()
	if err != nil {
		tx.Discard()
		return err
	}

	_, err = tx.Exec(ctx)
	if err != nil {
		tx.Discard()
		return err
	}

	return nil
}

func (store *UserStore) IsUserTokenAvaliable(token string) (bool, error) {
	ctx := context.Background()

	keys, err := store.rds.Keys(ctx, consts.REDIS_AVAILABLE_USER_TOKEN_LIST+":*").Result()
	if err != nil {
		return false, err
	}

	for _, key := range keys {

		tokens, err := store.rds.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return false, err
		}

		for _, t := range tokens {
			if t == token {
				return true, nil
			}
		}
	}

	return false, nil
}

func (store *UserStore) SaveUserAvatarByUID(uid uint64, fileName string, data []byte) error {
	savePath := filepath.Join(consts.AVATAR_IMAGE_PATH, fileName)

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return err
	}

	user := new(models.UserInfo)
	result := store.db.Where("id = ?", uid).First(user)
	if result.Error != nil {
		return result.Error
	}

	if user.Avatar != "vanilla.webp" {
		ctx := context.Background()
		store.rds.XAdd(ctx, &redis.XAddArgs{
			Stream: consts.AVATAR_CLEAN_STREAM,
			Values: map[string]interface{}{
				"filename": user.Avatar,
			},
		})
	}

	user.Avatar = fileName
	result = store.db.Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (store *UserStore) UpdateUserPasswordByUsername(username string, hashedNewPassword string) error {
	userAuthInfo := new(models.UserAuthInfo)
	result := store.db.Where("username = ?", username).First(userAuthInfo)
	if result.Error != nil {
		return result.Error
	}

	userAuthInfo.PasswordHash = hashedNewPassword
	result = store.db.Save(userAuthInfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (store *UserStore) UpdateUserInfoByUID(uid uint64, updatedProfile *models.UserInfo) error {
	var userProfile models.UserInfo
	result := store.db.Where("id = ?", uid).First(&userProfile)
	if result.Error != nil {
		return result.Error
	}

	userProfile.UpdatedAt = time.Now()
	userProfile.NickName = updatedProfile.NickName
	userProfile.Birth = updatedProfile.Birth
	userProfile.Gender = updatedProfile.Gender

	return store.db.Save(&userProfile).Error
}

func (store *UserStore) GetUserLikedRecord(uid int64) ([]int64, error) {
	postLikeCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_LIKE_COLLECTION)
	filter := bson.D{{Key: "uid", Value: uid}}
	sort := bson.D{{Key: "liked_at", Value: 1}}
	ctx := context.Background()
	defer ctx.Done()

	cursor, err := postLikeCollection.Find(ctx, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var postLikes []struct {
		PostID int64 `bson:"post_id"`
	}
	err = cursor.All(ctx, &postLikes)
	if err != nil {
		return nil, err
	}

	liked := make([]int64, len(postLikes))
	for index, postLike := range postLikes {
		liked[index] = postLike.PostID
	}

	return liked, nil
}

func (store *UserStore) GetUserFavoriteRecord(uid int64) ([]int64, error) {
	postFavoriteCollection := store.mongo.Database(consts.MONGODB_DATABASE_NAME).Collection(consts.POST_FAVORITE_COLLECTION)
	filter := bson.D{{Key: "uid", Value: uid}}
	sort := bson.D{{Key: "favourited_at", Value: 1}}
	ctx := context.Background()
	defer ctx.Done()

	cursor, err := postFavoriteCollection.Find(ctx, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var postFavorites []struct {
		PostID int64 `bson:"post_id"`
	}
	err = cursor.All(ctx, &postFavorites)
	if err != nil {
		return nil, err
	}

	favorited := make([]int64, len(postFavorites))
	for index, postFavorite := range postFavorites {
		favorited[index] = postFavorite.PostID
	}

	return favorited, nil
}
