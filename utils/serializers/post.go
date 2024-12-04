package serializers

import (
	"github.com/mehakhanaa/complex-micro-blog/models"
)

type PostListResponse struct {
	IDs []int64 `json:"ids"`
}

func NewPostListResponse(posts []int64) *PostListResponse {
	return &PostListResponse{IDs: posts}
}

type PostDetailResponse struct {
	CommentID    uint64   `json:"comment_id"`
	UID          uint64   `json:"uid"`
	Timestamp    int64    `json:"timestamp"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	ParentPostID *uint64  `json:"parent_post_id"`
	Images       []string `json:"images"`
	Like         int64    `json:"like"`
	Favourite    int64    `json:"favourite"`
	Farward      int      `json:"farward"`
}

func NewPostDetailResponse(post models.PostInfo, likeCount, favouriteCount int64) *PostDetailResponse {

	profileData := &PostDetailResponse{
		CommentID:    uint64(post.ID),
		UID:          post.UID,
		Timestamp:    post.CreatedAt.Unix(),
		Title:        post.Title,
		Content:      post.Content,
		ParentPostID: post.ParentPostID,
		Like:         likeCount,
		Favourite:    favouriteCount,
		Farward:      len(post.Farward),
	}
	for _, image := range post.Images {
		profileData.Images = append(profileData.Images, "/resources/image/"+image)
	}

	return profileData
}

type CreatePostResponse struct {
	ID uint64 `json:"id"`
}

func NewCreatePostResponse(postInfo models.PostInfo) CreatePostResponse {
	var resp = CreatePostResponse{
		ID: uint64(postInfo.ID),
	}
	return resp
}

type UploadPostImageResponse struct {
	UUID string `json:"uuid"`
}

func NewUploadPostImageResponse(uuid string) UploadPostImageResponse {
	return UploadPostImageResponse{UUID: uuid}
}

type PostUserStatus struct {
	PostID    uint64 `json:"post_id"`
	UID       uint64 `json:"uid"`
	Like      bool   `json:"like"`
	Favourite bool   `json:"favourite"`
}

func NewPostUserStatus(postID uint64, uid uint64, like bool, favourite bool) PostUserStatus {
	return PostUserStatus{
		PostID:    postID,
		UID:       uid,
		Like:      like,
		Favourite: favourite,
	}
}
