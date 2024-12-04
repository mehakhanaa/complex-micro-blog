package services

import (
	"context"
	"errors"
	"mime/multipart"
	"strconv"

	"github.com/lib/pq"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/models"
	search "github.com/mehakhanaa/complex-micro-blog/proto"
	"github.com/mehakhanaa/complex-micro-blog/stores"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/mehakhanaa/complex-micro-blog/utils/converters"
	"github.com/mehakhanaa/complex-micro-blog/utils/validers"
)

type PostService struct {
	postStore           *stores.PostStore
	searchServiceClient search.SearchEngineClient
}

func (factory *Factory) NewPostService(searchServiceClient search.SearchEngineClient) *PostService {
	return &PostService{
		postStore:           factory.storeFactory.NewPostStore(),
		searchServiceClient: searchServiceClient,
	}
}

func (service *PostService) GetPostList(reqType, uid, length, from string, userStore *stores.UserStore) ([]int64, error) {
	var (
		postInfos  []models.PostInfo
		userRecord pq.Int64Array
		err        error
		uidInt64   int64
	)

	if reqType != "all" {
		uidInt64, err = strconv.ParseInt(uid, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	queryLenth := 10
	if length != "" {
		queryLenth, err = strconv.Atoi(length)
		if err != nil {
			return nil, err
		}
		if queryLenth > 10 {
			queryLenth = 10
		}
	}

	switch reqType {
	case "all":
		postInfos, err = service.postStore.GetPostList(from, queryLenth)
	case "user":
		postInfos, err = service.postStore.GetPostListByUID(uid)
	case "liked":
		userRecord, err = userStore.GetUserLikedRecord(uidInt64)
	case "favourited":
		userRecord, err = userStore.GetUserFavoriteRecord(uidInt64)
	}
	if err != nil {
		return nil, err
	}

	if reqType == "all" || reqType == "user" {
		postIDs := make([]int64, len(postInfos))
		for index, post := range postInfos {
			postIDs[index] = int64(post.ID)
		}
		return postIDs, nil
	}

	if userRecord == nil {
		return nil, nil
	}
	postIDs := make([]int64, len(userRecord))
	for index, id := range userRecord {
		postIDs[index] = int64(id)
	}
	return postIDs, nil
}

func (service *PostService) GetPostInfo(postID uint64) (models.PostInfo, int64, int64, error) {
	return service.postStore.GetPostInfo(postID)
}

func (service *PostService) CreatePost(uid uint64, ipAddr string, postReqInfo types.PostCreateBody) (models.PostInfo, error) {

	for _, image := range postReqInfo.Images {
		existence, err := service.postStore.CheckCacheImageAvaliable(image)
		if err != nil {
			return models.PostInfo{}, err
		}
		if !existence {
			return models.PostInfo{}, errors.New("image does not exist")
		}
	}

	postInfo, err := service.postStore.CreatePost(uid, ipAddr, postReqInfo)
	if err != nil {
		return models.PostInfo{}, err
	}

	_, err = service.searchServiceClient.CreatePostIndex(context.TODO(), &search.CreatePostIndexRequest{
		Id:      int64(postInfo.ID),
		Title:   postReqInfo.Title,
		Content: postReqInfo.Content,
	})
	if err != nil {
		return models.PostInfo{}, err
	}

	return postInfo, nil
}

func (service *PostService) UploadPostImage(postImage *multipart.FileHeader) (string, error) {

	imageFile, err := postImage.Open()
	if err != nil {
		return "", err
	}
	defer imageFile.Close()

	fileType, err := validers.ValidImageFile(
		postImage,
		&imageFile,
		consts.POST_IMAGE_MIN_WIDTH,
		consts.POST_IMAGE_MIN_HEIGHT,
		consts.POST_IMAGE_MAX_FILE_SIZE,
	)
	if err != nil {
		return "", err
	}

	convertedImage, err := converters.ResizePostImage(fileType, &imageFile)
	if err != nil {
		return "", err
	}

	return service.postStore.CachePostImage(convertedImage)
}

func (service *PostService) LikePost(uid, postID int64) error {

	return service.postStore.LikePost(uid, postID)
}

func (service *PostService) CancelLikePost(uid, postID int64) error {

	return service.postStore.CancelLikePost(uid, postID)
}

func (service *PostService) FavouritePost(uid, postID int64) error {

	return service.postStore.FavouritePost(uid, postID)
}

func (service *PostService) CancelFavouritePost(uid, postID int64) error {

	return service.postStore.CancelFavouritePost(uid, postID)
}

func (service *PostService) GetPostUserStatus(uid, postID int64) (bool, bool, error) {

	return service.postStore.GetPostUserStatus(uid, postID)
}

func (service *PostService) DeletePost(postID uint64) error {

	return service.postStore.DeletePost(postID)
}
