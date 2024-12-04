package stores

import (
	search "github.com/mehakhanaa/complex-micro-blog/proto"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Factory struct {
	db            *gorm.DB
	rds           *redis.Client
	mongo         *mongo.Client
	searchService search.SearchEngineClient
}

func NewFactory(db *gorm.DB, redisClient *redis.Client, mongoClient *mongo.Client, searchService search.SearchEngineClient) *Factory {
	return &Factory{
		db:            db,
		rds:           redisClient,
		mongo:         mongoClient,
		searchService: searchService,
	}
}
