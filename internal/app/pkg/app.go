package app

import (
	"context"
	"os"

	"lab4/internal/app/config"
	"lab4/internal/app/dsn"
	redisClient "lab4/internal/app/redis"
	"lab4/internal/app/repository"
)

type Application struct {
	config   *config.Config
	repo     *repository.Repository
	redis    *redisClient.Client
	feedRepo repository.FeedRepository
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

	rds, err := redisClient.New(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}

	feedRepo := repository.NewFeedRepository(repo.DB(), rds.Raw())

	return &Application{
		config:   cfg,
		repo:     repo,
		redis:    rds,
		feedRepo: feedRepo,
	}, nil
}

func (a *Application) Run() error {
	a.StartServer()
	return nil
}