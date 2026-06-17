package app

import (
	"net/http"
    "context"
    "time"
    "fmt"
    "math/rand"
	"github.com/gin-gonic/gin"
)

func (a Application) GetFeed(c *gin.Context) {
	userIDAny, ok := c.Get("userid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDAny.(uint)

	ids, err := a.feedRepo.GetFeedForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(ids) == 0 {
		if err := a.feedRepo.GenerateFeedForUser(c.Request.Context(), userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ids, err = a.feedRepo.GetFeedForUser(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ids": ids,
	})
}

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