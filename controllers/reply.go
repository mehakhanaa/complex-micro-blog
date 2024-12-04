package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/services"
	"github.com/mehakhanaa/complex-micro-blog/stores"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/mehakhanaa/complex-micro-blog/utils/serializers"
)

type ReplyController struct {
	replyService *services.ReplyService
}

func (factory *Factory) NewReplyController() *ReplyController {
	return &ReplyController{
		replyService: factory.serviceFactory.NewReplyService(),
	}
}

func (controller *ReplyController) NewCreateReplyHandler(commentStore *stores.CommentStore, userStore *stores.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.ReplyCreateBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Content == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "content is required"),
			)
		}

		if reqBody.CommentID == 0 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err := controller.replyService.CreateReply(claims.UID, reqBody.CommentID, reqBody.ParentReplyID, reqBody.Content, commentStore, userStore)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "reply created successfully"),
		)
	}
}

func (controller *ReplyController) DeleteReplyHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserReplyDeleteBody)
		if err := ctx.BodyParser(reqBody); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.ReplyID == 0 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "reply id is required"),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		if err := controller.replyService.DeleteReply(claims.UID, reqBody.ReplyID); err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

func (controller *ReplyController) NewUpdateReplyHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserReplyUpdateBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Content == "" || reqBody.ReplyID == 0 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, " content or reply id is required"),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err = controller.replyService.UpdateReply(claims.UID, reqBody.ReplyID, reqBody.Content)
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

func (controller *ReplyController) NewGetReplyListHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commentID := ctx.Query("comment-id")
		if commentID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "comment id is required"),
			)
		}

		commentIDUint64, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		replyList, err := controller.replyService.GetReplyList(commentIDUint64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewReplyListResponse(replyList)),
		)
	}
}

func (controller *ReplyController) NewGetReplyDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		replyID := ctx.Query("reply-id")
		if replyID == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "reply id is required"),
			)
		}

		replyIDUint64, err := strconv.ParseUint(replyID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		reply, err := controller.replyService.GetReplyDetail(replyIDUint64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewReplyDetailResponse(reply)),
		)
	}
}
