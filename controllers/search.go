package controllers

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	search "github.com/mehakhanaa/complex-micro-blog/proto"
	"github.com/mehakhanaa/complex-micro-blog/services"
	"github.com/mehakhanaa/complex-micro-blog/utils/serializers"
)

type SearchController struct {
	searchService *services.SearchService
}

func (factory *Factory) NewSearchController(searchServiceClient search.SearchEngineClient) *SearchController {
	return &SearchController{
		searchService: factory.serviceFactory.NewSearchService(searchServiceClient),
	}
}

func (controller *SearchController) NewSearchPostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		queryString := ctx.Query("q")
		if queryString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "query content is required"),
			)
		}

		decodedQueryString, err := url.QueryUnescape(queryString)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "query content is invalid"),
			)
		}

		result, err := controller.searchService.SearchPost(decodedQueryString)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "", serializers.NewPostListResponse(result.Ids)),
		)
	}
}
