package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
	Minio       MinioConfig
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}           // создаем объект конфига
	err = viper.Unmarshal(cfg) // читаем информацию из файла,
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}

	// MinIO из env
	cfg.Minio.Endpoint = os.Getenv("MINIO_ENDPOINT")
	cfg.Minio.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	cfg.Minio.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	cfg.Minio.Bucket = os.Getenv("MINIO_BUCKET")
	cfg.Minio.UseSSL = os.Getenv("MINIO_USE_SSL") == "true"

	log.Info("config parsed")

	return cfg, nil
}
