package redis

import (
	"context"
	"encoding/json"
	"time"
	"lab4/internal/app/role"
)

const sessionPrefix = "session."

type Session struct {
	UserID uint      `json:"user_id"`
	Login  string    `json:"login"`
	Role   role.Role `json:"role"`
	ExpAt  time.Time `json:"exp_at"`
}

func getSessionKey(sessionID string) string {
	return servicePrefix + sessionPrefix + sessionID
}

func (c *Client) SaveSession(ctx context.Context, sessionID string, sess *Session, ttl time.Duration) error {
	data, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, getSessionKey(sessionID), data, ttl).Err()
}

func (c *Client) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	val, err := c.rdb.Get(ctx, getSessionKey(sessionID)).Result()
	if err != nil {
		return nil, err
	}

	var sess Session
	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		return nil, err
	}

	return &sess, nil
}

func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	return c.rdb.Del(ctx, getSessionKey(sessionID)).Err()
}
