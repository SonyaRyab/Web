package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
	Redis RedisConfig
	JWT JWTConfig
}

type JWTConfig struct {
	Token         string       
	ExpiresIn     time.Duration 
	SigningMethod string     
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
   envRedisHost = "REDIS_HOST"
   envRedisPort = "REDIS_PORT"
   envRedisUser = "REDIS_USER"
   envRedisPass = "REDIS_PASSWORD"
   envJWTToken = "JWT_TOKEN"
   envJWTExpires = "JWT_EXPIRES_HOURS"
)

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

	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.User = os.Getenv(envRedisUser)
	cfg.Redis.Password = os.Getenv(envRedisPass)

	cfg.JWT.Token = os.Getenv(envJWTToken)
	if cfg.JWT.Token == "" {
		cfg.JWT.Token = "default_secret_change_me" // дефолт для разработки
	}
	
	expiresHours := 24

	if os.Getenv(envRedisPort) != "" {
		cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))
		if err != nil {
			return nil, fmt.Errorf("redis port must be int value: %w", err)
		}
	}

	if os.Getenv(envJWTExpires) != "" {
		expiresHours, _ = strconv.Atoi(os.Getenv(envJWTExpires))
	}
	cfg.JWT.ExpiresIn = time.Duration(expiresHours) * time.Hour
	cfg.JWT.SigningMethod = "HS256"

	if cfg.Redis.Host == "" {
		cfg.Redis.Host = "127.0.0.1"
	}
	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}
	
	log.Info("config parsed")

	return cfg, nil
}
