package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/useragent"
	"gorm.io/gorm"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/models"
	"github.com/mehakhanaa/complex-micro-blog/services"
	"github.com/mehakhanaa/complex-micro-blog/types"
	"github.com/mehakhanaa/complex-micro-blog/utils/serializers"
)

type UserController struct {
	userService *services.UserService
}

func (factory *Factory) NewUserController() *UserController {
	return &UserController{
		userService: factory.serviceFactory.NewUserService(),
	}
}

func (controller *UserController) NewProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		uidString := ctx.Query("uid")
		username := ctx.Query("username")
		if uidString == "" && username == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "parameter uid or username is required"),
			)
		}

		var (
			user *models.UserInfo
			err  error
		)
		switch uidString {

		case "":
			user, err = controller.userService.GetUserInfoByUsername(username)

		default:
			var uid uint64

			if uid, err = strconv.ParseUint(uidString, 10, 64); err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
				)
			}
			user, err = controller.userService.GetUserInfoByUID(uid)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user does not exist"),
			)
		}

		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewUserProfileData(user)),
		)
	}
}

func (controller *UserController) NewRegisterHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserAuthBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Username == "" || reqBody.Password == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "username or password is required"),
			)
		}

		err = controller.userService.RegisterUser(reqBody.Username, reqBody.Password)
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

func (controller *UserController) NewLoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserAuthBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Username == "" || reqBody.Password == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "username or password is required"),
			)
		}

		userAgentString := ctx.Get("User-Agent")
		ua := useragent.New(userAgentString)

		browser, version := ua.Browser()
		var sb strings.Builder
		sb.WriteString(browser)
		sb.WriteString(" ")
		sb.WriteString(version)
		browserInfo := sb.String()

		os := ua.OSInfo().FullName

		token, err := controller.userService.LoginUser(reqBody.Username, reqBody.Password, ctx.IP(), browserInfo, os)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewUserToken(token)),
		)
	}
}

func (controller *UserController) NewUploadAvatarHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		form, err := ctx.MultipartForm()
		if err != nil {
			return err
		}
		files := form.File["avatar"]

		if len(files) != 1 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "required 1 file, but got more or less"),
			)
		}
		fileHeader := files[0]

		err = controller.userService.UserUploadAvatar(claims.UID, fileHeader)
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

func (controller *UserController) NewUpdatePasswordHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserUpdatePasswordBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.Username == "" || reqBody.Password == "" || reqBody.NewPassword == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "username, password or new password is required"),
			)
		}

		err = controller.userService.UserUpdatePassword(
			reqBody.Username,
			reqBody.Password,
			reqBody.NewPassword,
		)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "password updated successfully"),
		)
	}
}

func (controller *UserController) NewUpdateProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		reqBody := new(types.UserUpdateProfileBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		if reqBody.NickName == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "nickname is required"),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err = controller.userService.UpdateUserInfo(claims.UID, reqBody)
		if err != nil {
			return ctx.Status(500).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, "failed to update profile"),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "profile updated successfully"),
		)
	}
}
