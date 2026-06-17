package app

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func (a Application) BuildUserFeed(userID int64) error {
	ids, err := a.repo.GetAllReagentIDs()
	if err != nil {
		return err
	}

	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	if len(ids) > 100 {
		ids = ids[:100]
	}

	key := fmt.Sprintf("feed:user:%d", userID)

	if err := a.redis.Raw().Del(context.Background(), key).Err(); err != nil {
        return err
    }

	values := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		values = append(values, id)
	}

	if len(values) > 0 {
		if err := a.redis.Raw().RPush(context.Background(), key, values...).Err(); err != nil {
            return err
        }
		if err := a.redis.Raw().Expire(context.Background(), key, 24*time.Hour).Err(); err != nil {
            return err
        }
	}

	return nil
}