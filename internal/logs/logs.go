package logs

import (
	"context"
	"time"

	"github.com/n17ali/gohive/internal/storage"
)

type Logger interface {
	LogTaskExecution(ctx context.Context, taskID, status, message string) error
	GetTaskLogs(ctx context.Context, taskID string, count int64) ([]map[string]any, error)
}

type TaskLogger struct {
	store storage.LogStore
}

func NewTaskLogger(store storage.LogStore) Logger {
	return &TaskLogger{store: store}
}

func (l *TaskLogger) LogTaskExecution(ctx context.Context, taskID, status, message string) error {
	logEntry := map[string]any{
		"time":    time.Now().Format(time.RFC3339),
		"status":  status,
		"message": message,
	}
	return l.store.SaveLog(ctx, taskID, logEntry)
}

func (l *TaskLogger) GetTaskLogs(ctx context.Context, taskID string, count int64) ([]map[string]any, error) {
	return l.store.GetLogs(ctx, taskID, count)
}
