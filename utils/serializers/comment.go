package serializers

import (
	"github.com/mehakhanaa/complex-micro-blog/models"
)

type CommentListResponse struct {
	IDs []uint64 `json:"ids"`
}

func NewCommentListResponse(commentInfos []models.CommentInfo) CommentListResponse {
	var ids []uint64
	for _, commentInfos := range commentInfos {
		ids = append(ids, uint64(commentInfos.ID))
	}
	return CommentListResponse{IDs: ids}
}

type CommentDetailResponse struct {
	CommentID     uint64 `json:"comment_id"`
	PostID        uint64 `json:"post_id"`
	PosterUID     uint64 `json:"poster_uid"`
	PostTimestamp int64  `json:"post_timestamp"`
	Content       string `json:"content"`
	Likes         int64  `json:"likes"`
	Replies       int    `json:"replies"`
	Is_liked      bool   `json:"is_liked"`
	Is_disliked   bool   `json:"is_disliked"`
}

func NewCommentDetailResponse(comment models.CommentInfo, likeCount int64) *CommentDetailResponse {

	profileData := &CommentDetailResponse{
		CommentID:     uint64(comment.ID),
		PostID:        comment.PostID,
		PosterUID:     comment.UID,
		PostTimestamp: comment.CreatedAt.Unix(),
		Content:       comment.Content,
		Likes:         likeCount,
	}

	return profileData
}

type CreateCommentResponse struct {
	ID uint64 `json:"id"`
}

func NewCreateCommentResponse(commentID uint64) CreateCommentResponse {
	var resp = CreateCommentResponse{
		ID: commentID,
	}
	return resp
}

type CommentUserStatusResponse struct {
	IsLiked    bool `json:"is_liked"`
	IsDisliked bool `json:"is_disliked"`
}

func NewCommentUserStatusResponse(isLiked, isDisliked bool) CommentUserStatusResponse {
	return CommentUserStatusResponse{
		IsLiked:    isLiked,
		IsDisliked: isDisliked,
	}
}
