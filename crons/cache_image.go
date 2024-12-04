package crons

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type CachedImageCleanJob struct {
	logger *logrus.Logger
	rds    *redis.Client
}

func NewCachedImageCleanJob(logger *logrus.Logger, rds *redis.Client) *CachedImageCleanJob {
	return &CachedImageCleanJob{
		logger: logger,
		rds:    rds,
	}
}

func (job *CachedImageCleanJob) Run() {
	job.logger.Debugln("Caching Images init...")

	ctx := context.Background()

	keys, err := job.rds.Keys(ctx, consts.CACHE_IMAGE_LIST+":*").Result()
	if err != nil {
		job.logger.Errorln("Error in caching images:", err)
		return
	}

	for _, key := range keys {

		timestamp, err := job.rds.HGet(ctx, key, "expire").Int64()
		if err != nil {
			job.logger.Errorln("Error in caching images:", err)
			continue
		}

		filename, err := job.rds.HGet(ctx, key, "filename").Result()
		if err != nil {
			job.logger.Errorln("Error in caching images:", err)
			continue
		}

		if timestamp < time.Now().Unix() {
			tx := job.rds.TxPipeline()
			_, err := tx.XAdd(ctx, &redis.XAddArgs{
				Stream: consts.CACHE_IMG_CLEAN_STREAM,
				Values: map[string]interface{}{"filename": filename},
			}).Result()
			if err != nil {
				tx.Discard()
				job.logger.Errorln("Error in caching images:", err)
				return
			}

			_, err = tx.Del(ctx, key).Result()
			if err != nil {
				tx.Discard()
				job.logger.Errorln("Error in caching images:", err)
				return
			}

			_, err = tx.Exec(ctx)
			if err != nil {
				tx.Discard()
				job.logger.Errorln("Error in caching images:", err)
				continue
			}
		}
	}

	length, err := job.rds.XLen(ctx, consts.CACHE_IMG_CLEAN_STREAM).Result()
	if err != nil {
		job.logger.Errorln("Error in caching images:", err)
		return
	}
	if length == 0 {
		job.logger.Debugln("Error in length")
		return
	}

	messages, err := job.rds.XRead(ctx, &redis.XReadArgs{
		Streams: []string{consts.CACHE_IMG_CLEAN_STREAM, "0"},
		Count:   0,
		Block:   0,
	}).Result()
	if err != nil {
		job.logger.Errorln("Error in caching images:", err)
		return
	}

	for _, item := range messages[0].Messages {
		filename := item.Values["filename"].(string)
		err := os.Remove(filepath.Join(consts.POST_IMAGE_CACHE_PATH, filename))

		if errors.Is(err, os.ErrNotExist) {
			job.logger.Warningln("File doesnt exsist:", filename)
			_, err := job.rds.XDel(ctx, consts.CACHE_IMG_CLEAN_STREAM, item.ID).Result()
			if err != nil {
				job.logger.Errorln("Error in caching images:", err)
			}
			continue
		}
		if err != nil {
			job.logger.Warningln("Error in caching images:", err)
			continue
		}

		_, err = job.rds.XDel(ctx, consts.CACHE_IMG_CLEAN_STREAM, item.ID).Result()
		if err != nil {
			job.logger.Errorln("Error in caching images:", err)
		}
	}

	job.logger.Debugln("Caching Images Clean Job Done")
}
