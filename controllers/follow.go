package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/services"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/mehakhanaa/complex-micro-blog/utils/serializers"
)

type FollowController struct {
	followService *services.FollowService
}

func (factory *Factory) NewFollowController() *FollowController {
	return &FollowController{
		followService: factory.serviceFactory.NewFollowService(),
	}
}

func (controller *FollowController) NewCreateFollowHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		body := struct {
			UserID uint64 `json:"user_id" form:"user_id"`
		}{}
		err := ctx.BodyParser(&body)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"))
		}
		followedID := body.UserID

		if err := controller.followService.FollowUser(claims.UID, followedID); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

func (controller *FollowController) NewCancelFollowHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		body := struct {
			UserID uint64 `json:"user_id" form:"user_id"`
		}{}
		err := ctx.BodyParser(&body)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"))
		}
		followedID := body.UserID

		if err := controller.followService.CancelFollowUser(claims.UID, followedID); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

func (controller *FollowController) NewFollowListHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		follows, err := controller.followService.GetFollowList(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewFollowListResponse(follows)),
		)
	}
}

func (controller *FollowController) NewFollowCountHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		count, err := controller.followService.GetFollowCountByUID(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", struct{ Count int64 }{count}),
		)
	}
}

func (controller *FollowController) NewFollowerListHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		followers, err := controller.followService.GetFollowerList(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewFollowerListResponse(followers)),
		)
	}
}

func (controller *FollowController) NewFollowerCountHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userIDString := ctx.Query("user_id")
		if userIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is required"),
			)
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user_id is invalid"),
			)
		}

		count, err := controller.followService.GetFollowerCountByUID(userID)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", struct{ Count int64 }{count}),
		)
	}
}
