package types

type UserAuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdatePasswordBody struct {
	UserAuthBody
	NewPassword string `json:"new_password"`
}

type UserUpdateProfileBody struct {
	NickName *string `json:"nickname"`
	Birth    *uint64 `json:"birth"`
	Gender   *string `json:"gender"`
}

type UserCommentCreateBody struct {
	PostID  *uint64 `json:"post_id" form:"post_id"`
	Content string  `json:"content" form:"content"`
}

type UserCommentUpdateBody struct {
	CommentID *uint64 `json:"comment_id" form:"comment_id"`
	Content   string  `json:"content" form:"content"`
}

type UserPostInfo struct {
	UID   uint   `json:"id"`
	Title string `json:"title"`
}

type PostCreateBody struct {
	Title   string   `json:"title" form:"title"`
	Content string   `json:"content" form:"content"`
	Images  []string `json:"images" form:"images"`
}

type UserCommentDeleteBody struct {
	CommentID *uint64 `json:"comment_id" form:"comment_id"`
}

type ReplyCreateBody struct {
	CommentID     uint64 `json:"comment_id" form:"comment_id"`
	ParentReplyID uint64 `json:"parent_reply_id" form:"parent_reply_id"`
	Content       string `json:"content" form:"content"`
}

type UserReplyUpdateBody struct {
	ReplyID uint64 `json:"reply_id" form:"reply_id"`
	Content string `json:"content" form:"content"`
}

type UserReplyDeleteBody struct {
	ReplyID uint64 `json:"reply_id" form:"reply_id"`
}
