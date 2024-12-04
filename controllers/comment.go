package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/services"
	"github.com/mehakhanaa/complex-micro-blog/stores"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/mehakhanaa/complex-micro-blog/utils/serializers"
	"gorm.io/gorm"
)

type CommentController struct {
	commentService *services.CommentService
}

func (factory *Factory) NewCommentController() *CommentController {
	return &CommentController{
		commentService: factory.serviceFactory.NewCommentService(),
	}
}

func (controller *CommentController) NewCreateCommentHandler(postStore *stores.PostStore, userStore *stores.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserCommentCreateBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Content == "" || reqBody.PostID == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post_id or content is required"),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		commentID, err := controller.commentService.CreateComment(claims.UID, *reqBody.PostID, reqBody.Content, postStore, userStore)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "comment created successfully", serializers.NewCreateCommentResponse(commentID)),
		)
	}
}

func (controller *CommentController) NewUpdateCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserCommentUpdateBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Content == "" || reqBody.CommentID == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "content or comment id is required"),
			)
		}

		err = controller.commentService.UpdateComment(*reqBody.CommentID, reqBody.Content)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

func (controller *CommentController) DeleteCommentHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		reqBody := new(types.UserCommentDeleteBody)
		if err := c.BodyParser(reqBody); err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.CommentID == nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		if err := controller.commentService.DeleteComment(*reqBody.CommentID); err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return c.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

func (controller *CommentController) NewCommentListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		postID := c.Query("post-id")
		if postID == "" {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post id is required"),
			)
		}
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		comments, err := controller.commentService.GetCommentList(postIDUint)
		if err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return c.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewCommentListResponse(comments)),
		)
	}
}

func (controller *CommentController) NewCommentDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentIDString := ctx.Query("comment-id")
		if commentIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentID, err := strconv.ParseUint(commentIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		comment, likeCount, err := controller.commentService.GetCommentInfo(commentID)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment does not exist"),
			)
		}

		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewCommentDetailResponse(comment, likeCount)),
		)
	}
}

func (controller *CommentController) NewCommentUserStatusHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		liked, disliked, err := controller.commentService.GetCommentUserStatus(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewCommentUserStatusResponse(liked, disliked)),
		)
	}
}

func (controller *CommentController) NewLikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err = controller.commentService.LikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

func (controller *CommentController) NewCancelLikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err = controller.commentService.CancelLikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

func (controller *CommentController) NewDislikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err = controller.commentService.DislikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

func (controller *CommentController) NewCancelDislikeCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err = controller.commentService.CancelDislikeComment(claims.UID, commentIDUint)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}
