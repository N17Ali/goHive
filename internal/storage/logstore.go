package storage

import "context"

// LogStore provides structured logging operations
type LogStore interface {
	SaveLog(ctx context.Context, taskID string, logEntry map[string]interface{}) error
	GetLogs(ctx context.Context, taskID string, count int64) ([]map[string]interface{}, error)
}
