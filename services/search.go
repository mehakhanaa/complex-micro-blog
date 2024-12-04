package services

import (
	"context"

	search "github.com/mehakhanaa/complex-micro-blog/proto"
)

type SearchService struct {
	searchServiceClient search.SearchEngineClient
}

func (factory *Factory) NewSearchService(searchServiceClient search.SearchEngineClient) *SearchService {
	return &SearchService{
		searchServiceClient: searchServiceClient,
	}
}

func (service *SearchService) SearchPost(queryString string) (*search.SearchResponse, error) {
	return service.searchServiceClient.Search(context.TODO(), &search.SearchRequest{
		Query: queryString,
	})
}
