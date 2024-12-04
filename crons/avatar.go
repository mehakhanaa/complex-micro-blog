package crons

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AvatarCleanJob struct {
	logger *logrus.Logger
	rds    *redis.Client
}

func NewAvatarCleanJob(logger *logrus.Logger, redisClient *redis.Client) *AvatarCleanJob {
	return &AvatarCleanJob{
		logger: logger,
		rds:    redisClient,
	}
}

func (job *AvatarCleanJob) Run() {
	job.logger.Debugln("Avatar clean job init...")

	ctx := context.Background()

	length, err := job.rds.XLen(ctx, consts.AVATAR_CLEAN_STREAM).Result()
	if err != nil {
		job.logger.Errorln("Error in avart clean job:", err)
		return
	}
	if length == 0 {
		job.logger.Debugln("Error in avart clean job")
		return
	}

	messages, err := job.rds.XRead(ctx, &redis.XReadArgs{
		Streams: []string{consts.AVATAR_CLEAN_STREAM, "0"},
		Count:   0,
		Block:   0,
	}).Result()
	if err != nil {
		job.logger.Errorln("Error in avart clean job:", err)
		return
	}

	for _, item := range messages[0].Messages {
		filename := item.Values["filename"].(string)
		err := os.Remove(filepath.Join(consts.AVATAR_IMAGE_PATH, filename))

		if errors.Is(err, os.ErrNotExist) {
			job.logger.Warnln("File doesnt exists:", filename)
			_, err := job.rds.XDel(ctx, consts.AVATAR_CLEAN_STREAM, item.ID).Result()
			if err != nil {
				job.logger.Errorln("Error in avart clean job:", err)
			}
			continue
		}
		if err != nil {
			job.logger.Warningln("Error in avart clean job:", err)
			continue
		}

		_, err = job.rds.XDel(ctx, consts.AVATAR_CLEAN_STREAM, item.ID).Result()
		if err != nil {
			job.logger.Errorln("Error in avart clean job:", err)
			continue
		}
	}

	job.logger.Debugln("Avatar clean job done")
}
