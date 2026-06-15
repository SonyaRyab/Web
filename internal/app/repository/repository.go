package repository

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

    // "gorm.io/gorm/logger"
	// "log"
    // "os"
    // "time"
)

type Repository struct {
	db                *gorm.DB
	minio             *minio.Client
	minio_bucket_name string
}

type RepositorySettings struct {
	PostgresDSN     string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

func New(settings *RepositorySettings) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(settings.PostgresDSN), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	useSSL := false // при true подключаемся к MinIO по HTTPS

	minioClient, err := minio.New(settings.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.MinioAccessKey, settings.MinioSecretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	// Возвращаем указатель на получившийся Repository
	return &Repository{
		db:                db,
		minio:             minioClient,
		minio_bucket_name: settings.MinioBucketName,
	}, nil
}

// func New(settings *RepositorySettings) (*Repository, error) {
//     newLogger := logger.New(
//         log.New(os.Stdout, "\r\n", log.LstdFlags),
//         logger.Config{
//             SlowThreshold:             time.Second,
//             LogLevel:                  logger.Info,
//             IgnoreRecordNotFoundError: true,
//             Colorful:                  false,
//         },
//     )

//     db, err := gorm.Open(postgres.Open(settings.PostgresDSN), &gorm.Config{
//         Logger: newLogger,
//     })
//     if err != nil {
//         return nil, err
//     }

//     useSSL := false

//     minioClient, err := minio.New(settings.MinioEndpoint, &minio.Options{
//         Creds:  credentials.NewStaticV4(settings.MinioAccessKey, settings.MinioSecretKey, ""),
//         Secure: useSSL,
//     })
//     if err != nil {
//         return nil, err
//     }

//     return &Repository{
//         db:                db,
//         minio:             minioClient,
//         minio_bucket_name: settings.MinioBucketName,
//     }, nil
// }

func (r *Repository) DB() *gorm.DB {
    return r.db
}