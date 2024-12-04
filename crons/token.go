package crons

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/utils/parsers"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type TokenCleanJob struct {
	logger *logrus.Logger
	rds    *redis.Client
}

func NewTokenCleanJob(logger *logrus.Logger, redisClient *redis.Client) *TokenCleanJob {
	return &TokenCleanJob{
		logger: logger,
		rds:    redisClient,
	}
}

func (job *TokenCleanJob) Run() {
	job.logger.Debugln("Token clean job init......")

	ctx := context.Background()

	keys, err := job.rds.Keys(ctx, consts.REDIS_AVAILABLE_USER_TOKEN_LIST+":*").Result()
	if err != nil {
		job.logger.Errorln("Error in token clean job:", err)
		return
	}

	for _, key := range keys {

		tokens, err := job.rds.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			job.logger.Errorln("Error in token clean job:", err)
			continue
		}

		for _, token := range tokens {
			_, err := parsers.ParseToken(token)

			if errors.Is(err, jwt.ErrTokenExpired) {

				_, err := job.rds.LRem(ctx, key, 0, token).Result()
				if err != nil {
					job.logger.Errorln("Error in token clean job:", err)
				}
				continue
			}

			if err != nil {
				job.logger.Errorln("Error in token clean job:", err)
				continue
			}
		}
	}

	job.logger.Debugln("Token clean job done")
}
