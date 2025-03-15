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
	return redis.Client.Set(ctx, "task:"+task.ID, data, 0).Err()
}

func GetTask(ctx context.Context, id string) (*Task, error) {
	data, err := redis.Client.Get(ctx, "task:"+id).Result()
	if err != nil {
		return nil, err
	}
	var task Task
	if err = json.Unmarshal([]byte(data), &task); err != nil {
		return nil, err
	}
	task.Interval = time.Duration(task.Interval) * time.Second

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
	if err := redis.Client.Del(ctx, "task:"+id).Err(); err != nil {
		return err
	}
	return nil
}

func GetAllTask(ctx context.Context) ([]Task, error) {
	var tasks []Task

	keys, err := redis.Client.Keys(ctx, "task:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		taskID := key[len("task:"):]

		task, err := GetTask(ctx, taskID)
		if err == nil {
			tasks = append(tasks, *task)
		}
	}

	return tasks, nil
}

func GetTaskLastRunTime(ctx context.Context, taskID string) (time.Time, error) {
	ts, err := redis.Client.Get(ctx, "lastrun:"+taskID).Int64()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(ts, 0), nil
}

func SetTaskLatRunTime(ctx context.Context, taskID string, t time.Time) error {
	return redis.Client.Set(ctx, "lastrun:"+taskID, t.Unix(), 0).Err()
}
