package repository

import (
    "context"
    "fmt"
    "math/rand"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"

    "lab4/internal/app/ds"
)

type FeedRepository struct {
    db    *gorm.DB
    redis *redis.Client
}

func NewFeedRepository(db *gorm.DB, redisClient *redis.Client) *FeedRepository {
    return &FeedRepository{db: db, redis: redisClient}
}

// GenerateFeedForUser формирует 100 случайных id реагентов и сохраняет в Redis.
func (r *FeedRepository) GenerateFeedForUser(ctx context.Context, userID uint) error {
    var reagents []ds.Reagent
    if err := r.db.Select("id").Find(&reagents).Error; err != nil {
        return err
    }
    if len(reagents) == 0 {
        return nil
    }

    rand.Seed(time.Now().UnixNano())
    // перемешиваем
    rand.Shuffle(len(reagents), func(i, j int) {
        reagents[i], reagents[j] = reagents[j], reagents[i]
    })

    // берём первые 100
    count := 100
    if len(reagents) < count {
        count = len(reagents)
    }

    key := feedKey(userID)
    // чистим старую ленту
    if err := r.redis.Del(ctx, key).Err(); err != nil {
        return err
    }

    // добавляем как ZSET: score = позиция, member = id
    members := make([]*redis.Z, 0, count)
    for i := 0; i < count; i++ {
        members = append(members, &redis.Z{
            Score:  float64(i),
            Member: reagents[i].ID,
        })
    }

    if err := r.redis.ZAdd(ctx, key, members...).Err(); err != nil {
        return err
    }

    // опционально срок жизни ленты
    _ = r.redis.Expire(ctx, key, 24*time.Hour).Err()

    return nil
}

// GetFeedForUser возвращает весь массив id из Redis.
func (r *FeedRepository) GetFeedForUser(ctx context.Context, userID uint) ([]uint, error) {
    key := feedKey(userID)
    values, err := r.redis.ZRange(ctx, key, 0, -1).Result()
    if err != nil {
        return nil, err
    }

    ids := make([]uint, 0, len(values))
    for _, v := range values {
        id64, err := strconv.ParseUint(v, 10, 64)
        if err != nil {
            continue
        }
        ids = append(ids, uint(id64))
    }
    return ids, nil
}

func feedKey(userID uint) string {
    return fmt.Sprintf("feed:%d", userID)
}