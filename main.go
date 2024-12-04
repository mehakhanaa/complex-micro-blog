package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/mehakhanaa/complex-micro-blog/configs"
	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/controllers"
	"github.com/mehakhanaa/complex-micro-blog/crons"
	"github.com/mehakhanaa/complex-micro-blog/loggers"
	"github.com/mehakhanaa/complex-micro-blog/middlewares"
	"github.com/mehakhanaa/complex-micro-blog/models"
	search "github.com/mehakhanaa/complex-micro-blog/proto"
	"github.com/mehakhanaa/complex-micro-blog/services"
	"github.com/mehakhanaa/complex-micro-blog/stores"
)

var (
	logger              *logrus.Logger
	cfg                 *configs.Config
	db                  *gorm.DB
	redisClient         *redis.Client
	mongoClient         *mongo.Client
	searchSeviceConn    *grpc.ClientConn
	searchServiceClient search.SearchEngineClient
	storeFactory        *stores.Factory
	controllerFactory   *controllers.Factory
	middlewareFactory   *middlewares.Factory
)

func init() {

	logger = loggers.NewLogger()
	logger.Infoln("Starting.....")

	var err error

	cfg, err = configs.NewConfig()
	if err != nil {
		logger.Panicln(err.Error())
	}

	var (
		logLevel logrus.Level
		logMode  gormLogger.LogLevel
	)
	switch cfg.Env.Type {
	case "development":
		logLevel = logrus.DebugLevel
		logMode = gormLogger.Error
	case "production":
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	default:
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	}

	logger.SetLevel(logLevel)
	logger.Debugln("Connecting to db:", strings.ToUpper(logLevel.String()))

	db, err = gorm.Open(
		postgres.Open(fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(logMode),
		},
	)
	if err != nil {
		logger.Panicln(err.Error())
	}

	logger.Debugln("Migrating DB...")
	err = models.Migrate(db)
	if err != nil {
		logger.Panicln("Migrating error", err.Error())
	}

	logger.Debugln("Init Redis...")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	_, err = redisClient.Ping(context.TODO()).Result()
	if err != nil {
		logger.Panicln("Error in Redis: ", err.Error())
	}
	logger.Debugln("Redis Connected")

	logger.Debugln("Init MongoDB...")
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logger.Panicln("Error MongoDB:", err.Error())
	}
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		logger.Panicln("Error in ping MongoDB", err.Error())
	}
	logger.Debugln("MongoDB Connected")

	searchSeviceConn, err = grpc.Dial(fmt.Sprintf("%s:%d", cfg.SearchService.Host, cfg.SearchService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Panicln("grpc error", err.Error())
	}
	searchServiceClient = search.NewSearchEngineClient(searchSeviceConn)

	storeFactory = stores.NewFactory(db, redisClient, mongoClient, searchServiceClient)

	controllerFactory = controllers.NewFactory(
		services.NewFactory(storeFactory),
	)

	middlewareFactory = middlewares.NewFactory(storeFactory)
}

func main() {

	crons.InitJobs(logger, db, redisClient)

	var fiberConfig fiber.Config

	if cfg.Env.Type == "production" {
		fiberConfig = fiber.Config{
			Prefork: true,
		}
	}
	fiberConfig.BodyLimit = consts.REQUEST_BODY_LIMIT
	app := fiber.New(fiberConfig)

	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}][${latency}][${status}][${method}] ${path}\n",
	}))
	app.Use(compress.New(compress.Config{
		Level: cfg.Compress.Level,
	}))

	authMiddleware := middlewareFactory.NewTokenAuthMiddleware()

	resource := app.Group("/resources")

	resource.Static("/avatar", consts.AVATAR_IMAGE_PATH, fiber.Static{
		Compress: true,
	})

	resource.Static("/image", consts.POST_IMAGE_PATH, fiber.Static{
		Compress: true,
	})

	api := app.Group("/api")

	userController := controllerFactory.NewUserController()
	user := api.Group("/user")
	user.Get("/profile", userController.NewProfileHandler())
	user.Post("/register", userController.NewRegisterHandler())
	user.Post("/login", userController.NewLoginHandler())
	user.Post("/upload-avatar", authMiddleware.NewMiddleware(), userController.NewUploadAvatarHandler())
	user.Post("/update-psw", userController.NewUpdatePasswordHandler())
	user.Post("/edit", authMiddleware.NewMiddleware(), userController.NewUpdateProfileHandler())

	postController := controllerFactory.NewPostController(searchServiceClient)
	post := api.Group("/post")
	post.Get("/list", postController.NewPostListHandler(storeFactory.NewUserStore()))
	post.Get("/user-status", authMiddleware.NewMiddleware(), postController.NewPostUserStatusHandler())
	post.Post("/new", authMiddleware.NewMiddleware(), postController.NewCreatePostHandler())
	post.Post("/upload-img", authMiddleware.NewMiddleware(), postController.NewUploadPostImageHandler())
	post.Post("/like", authMiddleware.NewMiddleware(), postController.NewLikePostHandler())
	post.Post("/cancel-like", authMiddleware.NewMiddleware(), postController.NewCancelLikePostHandler())
	post.Post("/favourite", authMiddleware.NewMiddleware(), postController.NewFavouritePostHandler())
	post.Post("/cancel-favourite", authMiddleware.NewMiddleware(), postController.NewCancelFavouritePostHandler())
	post.Get("/:post", postController.NewPostDetailHandler())
	post.Delete("/:post", authMiddleware.NewMiddleware(), postController.NewDeletePostHandler())

	commentController := controllerFactory.NewCommentController()
	comment := api.Group("/comment")
	comment.Get("/list", commentController.NewCommentListHandler())
	comment.Get("/detail", commentController.NewCommentDetailHandler())
	comment.Get("/user-status", authMiddleware.NewMiddleware(), commentController.NewCommentUserStatusHandler())
	comment.Post("/edit", authMiddleware.NewMiddleware(), commentController.NewUpdateCommentHandler())
	comment.Post("/delete", authMiddleware.NewMiddleware(), commentController.DeleteCommentHandler())
	comment.Post("/like", authMiddleware.NewMiddleware(), commentController.NewLikeCommentHandler())
	comment.Post("/cancel-like", authMiddleware.NewMiddleware(), commentController.NewCancelLikeCommentHandler())
	comment.Post("/dislike", authMiddleware.NewMiddleware(), commentController.NewDislikeCommentHandler())
	comment.Post("/cancel-dislike", authMiddleware.NewMiddleware(), commentController.NewCancelDislikeCommentHandler())
	comment.Post("/new", authMiddleware.NewMiddleware(), commentController.NewCreateCommentHandler(
		storeFactory.NewPostStore(),
		storeFactory.NewUserStore(),
	))

	replyController := controllerFactory.NewReplyController()
	reply := api.Group("/reply")
	reply.Get("/list", replyController.NewGetReplyListHandler())
	reply.Get("/detail", replyController.NewGetReplyDetailHandler())
	reply.Post("/new", authMiddleware.NewMiddleware(), replyController.NewCreateReplyHandler(
		storeFactory.NewCommentStore(),
		storeFactory.NewUserStore()),
	)
	reply.Post("/edit", authMiddleware.NewMiddleware(), replyController.NewUpdateReplyHandler())
	reply.Post("/delete", authMiddleware.NewMiddleware(), replyController.DeleteReplyHandler())

	searchController := controllerFactory.NewSearchController(searchServiceClient)
	search := api.Group("/search")
	search.Get("/post", searchController.NewSearchPostHandler())

	followController := controllerFactory.NewFollowController()
	follow := api.Group("/follow")
	follow.Post("/new", authMiddleware.NewMiddleware(), followController.NewCreateFollowHandler())
	follow.Post("/delete", authMiddleware.NewMiddleware(), followController.NewCancelFollowHandler())
	follow.Get("/list", followController.NewFollowListHandler())
	follow.Get("/list-count", followController.NewFollowCountHandler())
	follow.Get("/follower-list", followController.NewFollowerListHandler())
	follow.Get("/follower-list-count", followController.NewFollowerCountHandler())

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Server.Port)))
}
