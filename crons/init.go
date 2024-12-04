package crons

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/mehakhanaa/complex-micro-blog/utils/jobs"
)

func InitJobs(logger *logrus.Logger, db *gorm.DB, redisClient *redis.Client) {

	crontab := cron.New()

	_, err := jobs.AddSkipIfStillRunningJob(crontab, "@every 5m", NewAvatarCleanJob(logger, redisClient))
	if err != nil {
		logger.Panicln(err.Error())
	}

	_, err = jobs.AddSkipIfStillRunningJob(crontab, "@every 5m", NewCachedImageCleanJob(logger, redisClient))
	if err != nil {
		logger.Panicln(err.Error())
	}

	_, err = jobs.AddSkipIfStillRunningJob(crontab, "@every 1h", NewTokenCleanJob(logger, redisClient))
	if err != nil {
		logger.Panicln(err.Error())
	}

	crontab.Start()
}
