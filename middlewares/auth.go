package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/stores"
	"github.com/mehakhanaa/complex-micro-blog/utils/parsers"
	"github.com/mehakhanaa/complex-micro-blog/utils/serializers"
)

type TokenAuthMiddleware struct {
	userStore *stores.UserStore
}

func (factory *Factory) NewTokenAuthMiddleware() *TokenAuthMiddleware {
	return &TokenAuthMiddleware{userStore: factory.store.NewUserStore()}
}

func (middleware *TokenAuthMiddleware) NewMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		token := ctx.Get("Authorization")
		if token == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "bearer token is required"),
			)
		}
		if len(token) < 7 || token[:7] != "Bearer " {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "bearer token is invalid"),
			)
		}
		token = token[7:]

		claims, err := parsers.ParseToken(token)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, "bearer token is expired"),
			)
		}
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, err.Error()),
			)
		}

		isAvaliable, err := middleware.userStore.IsUserTokenAvaliable(token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		if !isAvaliable {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, "bearer token is not avaliable"),
			)
		}

		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}
