package redis

import (
    "context"
    "fmt"
    "strconv"

    "lab4/internal/app/config"
    "github.com/go-redis/redis/v8"
)

const servicePrefix = "lab4_service."

type Client struct {
    cfg config.RedisConfig
    rdb *redis.Client
}

func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Password:   cfg.Password,
        Username:   cfg.User,
        Addr:       cfg.Host + ":" + strconv.Itoa(cfg.Port),
        DB:         0,
        DialTimeout: cfg.DialTimeout,
        ReadTimeout: cfg.ReadTimeout,
    })

    if _, err := rdb.Ping(ctx).Result(); err != nil {
        return nil, fmt.Errorf("cant ping redis: %w", err)
    }

    return &Client{
        cfg: cfg,
        rdb: rdb,
    }, nil
}

func (c *Client) Close() error {
    return c.rdb.Close()
}

func (c *Client) Raw() *redis.Client {
    return c.rdb
}