package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	search "github.com/mehakhanaa/complex-micro-blog/proto"
	"github.com/mehakhanaa/complex-micro-blog/services"
	"github.com/mehakhanaa/complex-micro-blog/stores"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/mehakhanaa/complex-micro-blog/utils/functools"
	"github.com/mehakhanaa/complex-micro-blog/utils/serializers"
)

type PostController struct {
	postService *services.PostService
}

func (factory *Factory) NewPostController(searchServiceClient search.SearchEngineClient) *PostController {
	return &PostController{
		postService: factory.serviceFactory.NewPostService(searchServiceClient),
	}
}

func (controller *PostController) NewPostListHandler(userStore *stores.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqType := ctx.Query("type")
		uid := ctx.Query("uid")
		length := ctx.Query("len")
		from := ctx.Query("from")
		if reqType == "user" || reqType == "liked" || reqType == "favourited" {
			_, err := strconv.ParseUint(uid, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid uid"),
				)
			}
		}
		if length != "" {
			_, err := strconv.ParseUint(length, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid length"),
				)
			}
		}
		if from != "" {
			_, err := strconv.ParseUint(from, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "invalid from id"),
				)
			}
		}

		var (
			posts []int64
			err   error
		)
		switch reqType {
		case "":
			posts, err = controller.postService.GetPostList("all", "", length, from, userStore)
		case "all":
			posts, err = controller.postService.GetPostList("all", "", length, from, userStore)
		case "user":
			posts, err = controller.postService.GetPostList("user", uid, length, from, userStore)
		case "liked":
			posts, err = controller.postService.GetPostList("liked", uid, length, from, userStore)
			posts = functools.Reverse(posts)
		case "favourited":
			posts, err = controller.postService.GetPostList("favourited", uid, length, from, userStore)
			posts = functools.Reverse(posts)
		default:
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "invalid type"))
		}
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewPostListResponse(posts)),
		)
	}
}

func (controller *PostController) NewPostDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		postIDString := ctx.Params("post")
		if postIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post id is required"),
			)
		}

		postID, err := strconv.ParseUint(postIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		post, likeCount, favouriteCount, err := controller.postService.GetPostInfo(postID)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post does not exist"),
			)
		}

		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewPostDetailResponse(post, likeCount, favouriteCount)),
		)
	}
}

func (controller *PostController) NewCreatePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		reqBody := types.PostCreateBody{}
		err := ctx.BodyParser(&reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Title == "" || reqBody.Content == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post title or post content is required"),
			)
		}
		if len(reqBody.Images) > 9 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post images count exceeds the limit"),
			)
		}

		postInfo, err := controller.postService.CreatePost(claims.UID, ctx.IP(), reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(
				consts.SUCCESS,
				"post created successfully",
				serializers.NewCreatePostResponse(postInfo),
			),
		)
	}
}

func (controller *PostController) NewPostUserStatusHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		postID := ctx.Query("post-id")

		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		isLiked, isFavourited, err := controller.postService.GetPostUserStatus(int64(claims.UID), int64(postIDUint))
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.Status(200).JSON(serializers.NewResponse(
			consts.SUCCESS,
			"succeed",
			serializers.NewPostUserStatus(
				postIDUint,
				claims.UID,
				isLiked,
				isFavourited,
			),
		),
		)
	}
}

func (controller *PostController) NewUploadPostImageHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		files := form.File["file"]
		if len(files) < 1 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "image is required"),
			)
		}
		if len(files) > 1 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "the number of image cannot exceed 1"),
			)
		}

		UUID, err := controller.postService.UploadPostImage(files[0])
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(serializers.NewResponse(
			consts.SUCCESS,
			"succeed",
			serializers.NewUploadPostImageResponse(UUID),
		))
	}
}

func (controller *PostController) NewLikePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		postID := ctx.Query("post-id")

		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		if err := controller.postService.LikePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

func (controller *PostController) NewCancelLikePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		postID := ctx.Query("post-id")

		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		if err := controller.postService.CancelLikePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

func (controller *PostController) NewFavouritePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		postID := ctx.Query("post-id")

		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		if err := controller.postService.FavouritePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

func (controller *PostController) NewCancelFavouritePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		postID := ctx.Query("post-id")

		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		if err := controller.postService.CancelFavouritePost(int64(claims.UID), int64(postIDUint)); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}

func (controller *PostController) NewDeletePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		postID := ctx.Params("post")

		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		if err := controller.postService.DeletePost(postIDUint); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}
