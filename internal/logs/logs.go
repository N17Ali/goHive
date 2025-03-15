package logs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/n17ali/gohive/pkg/redis"
)

type Logger interface {
	LogTaskExecution(ctx context.Context, taskID string, status string, message string) error
	GetTaskLogs(ctx context.Context, taskID string, count int64) ([]map[string]string, error)
}

type RedisLogger struct{}

func (r RedisLogger) LogTaskExecution(ctx context.Context, taskID, status, message string) error {
	logEntiry := map[string]any{
		"time":    time.Now().Format(time.RFC3339),
		"status":  status,
		"message": message,
	}

	logData, err := json.Marshal(logEntiry)
	if err != nil {
		return err
	}

	key := "log:" + taskID
	return redis.Client.LPush(ctx, key, logData).Err()
}

func (r RedisLogger) GetTaskLogs(ctx context.Context, taskID string, count int64) ([]map[string]string, error) {
	key := "log:" + taskID

	logs, err := redis.Client.LRange(ctx, key, 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	var results []map[string]string
	for _, logEntriy := range logs {
		var logData map[string]string
		if err = json.Unmarshal([]byte(logEntriy), &logData); err == nil {
			results = append(results, logData)
		}
	}
	return results, nil
}
