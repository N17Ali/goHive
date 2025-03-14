package task

import (
	context "context"
	"encoding/json"
	"time"

	"github.com/n17ali/gohive/pkg/redis"
)

type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Interval    time.Duration `json:"inteval"`
}

func SaveTask(ctx context.Context, task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return redis.Client.Set(ctx, task.ID, data, 0).Err()
}

func GetTask(ctx context.Context, id string) (*Task, error) {
	data, err := redis.Client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	var task Task
	if err = json.Unmarshal([]byte(data), &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(ctx context.Context, id string, task Task) error {
	existingTask, err := GetTask(ctx, id)
	if err != nil {
		return err
	}

	existingTask.Title = task.Title
	existingTask.Description = task.Description
	existingTask.Interval = task.Interval

	return SaveTask(ctx, *existingTask)
}

func DeleteTask(ctx context.Context, id string) error {
	if err := redis.Client.Del(ctx, id).Err(); err != nil {
		return err
	}
	return nil
}
