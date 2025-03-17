package redis

import (
	"context"
	"encoding/json"

	"github.com/n17ali/gohive/internal/storage"
	"github.com/redis/go-redis/v9"
)

type RedisLogStore struct {
	client *redis.Client
}

func NewRedisLogStore(client *redis.Client) storage.LogStore {
	return &RedisLogStore{client: client}
}

func (r *RedisLogStore) SaveLog(ctx context.Context, taskID string, logEntry map[string]interface{}) error {
	data, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, "log:"+taskID, data).Err()
}

func (r *RedisLogStore) GetLogs(ctx context.Context, taskID string, count int64) ([]map[string]interface{}, error) {
	rawLogs, err := r.client.LRange(ctx, "log:"+taskID, 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	var logs []map[string]interface{}
	for _, rawLog := range rawLogs {
		var logEntry map[string]interface{}
		if err = json.Unmarshal([]byte(rawLog), &logEntry); err == nil {
			logs = append(logs, logEntry)
		}
	}
	return logs, nil
}
