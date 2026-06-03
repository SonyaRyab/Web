package app

import (
	"context"
	"log"
	"os"
	"lab4/internal/app/config"
	"lab4/internal/app/dsn"
	"lab4/internal/app/redis"
	"lab4/internal/app/repository"
)

type Application struct {
	config *config.Config
	repo   *repository.Repository
	redis  *redis.Client
	feedRepo *repository.FeedRepository
}

func New(ctx context.Context) (*Application, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(&repository.RepositorySettings{
		PostgresDSN:     dsn.FromEnv(),
		MinioEndpoint:   os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		MinioBucketName: os.Getenv("MINIO_BUCKET"),
	})
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}
	feedRepo := repository.NewFeedRepository(repo.DB(), redisClient.Raw())

	return &Application{
		config: cfg,
		repo:   repo,
		redis: redisClient,
		feedRepo: feedRepo,
	}, nil
}

func (a *Application) Run() error {
	log.Println("application start running")
	a.StartServer()
	log.Println("application shut down")
	return nil
}
