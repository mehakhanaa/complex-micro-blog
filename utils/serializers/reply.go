package serializers

import "github.com/mehakhanaa/complex-micro-blog/models"

type ReplyListResponse struct {
	IDs []uint64 `json:"ids"`
}

func NewReplyListResponse(replies []uint64) ReplyListResponse {
	return ReplyListResponse{IDs: replies}
}

type ReplyDetailResponse struct {
	CreateTime     int64   `json:"create_time"`
	CommentID      uint64  `json:"comment_id"`
	UID            uint64  `json:"uid"`
	ParentReplyID  *uint64 `json:"parent_reply_id"`
	ParentReplyUID *uint64 `json:"parent_reply_uid"`
	Content        string  `json:"content"`
}

func NewReplyDetailResponse(reply models.ReplyInfo) ReplyDetailResponse {

	profileData := ReplyDetailResponse{
		CreateTime:     reply.CreatedAt.Unix(),
		CommentID:      reply.CommentID,
		UID:            reply.UID,
		ParentReplyID:  reply.ParentReplyID,
		ParentReplyUID: reply.ParentReplyUID,
		Content:        reply.Content,
	}

	return profileData
}
