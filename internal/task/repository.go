package task

import (
	context "context"
	"encoding/json"
	"time"

	"github.com/n17ali/gohive/internal/storage"
)

type TaskRepository struct {
	store storage.Store
}

func NewTaskRepository(store storage.Store) *TaskRepository {
	return &TaskRepository{
		store: store,
	}
}

func (r *TaskRepository) SaveTask(ctx context.Context, task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.store.Set(ctx, "task:"+task.ID, data)
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*Task, error) {
	data, err := r.store.Get(ctx, "task:"+id)
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

func (r *TaskRepository) UpdateTask(ctx context.Context, id string, task Task) error {
	existingTask, err := r.GetTask(ctx, id)
	if err != nil {
		return err
	}

	existingTask.Title = task.Title
	existingTask.Description = task.Description
	existingTask.Interval = task.Interval

	return r.SaveTask(ctx, *existingTask)
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	if err := r.store.Del(ctx, "task:"+id); err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) GetAllTask(ctx context.Context) ([]Task, error) {
	var tasks []Task

	keys, err := r.store.Keys(ctx, "task:*")
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		taskID := key[len("task:"):]

		task, err := r.GetTask(ctx, taskID)
		if err == nil {
			tasks = append(tasks, *task)
		}
	}

	return tasks, nil
}

func (r *TaskRepository) GetTaskLastRunTime(ctx context.Context, taskID string) (time.Time, error) {
	ts, err := r.store.Get(ctx, "lastrun:"+taskID)
	if err != nil {
		return time.Time{}, err
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", ts)
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

func (r *TaskRepository) SetTaskLatRunTime(ctx context.Context, taskID string, t time.Time) error {
	return r.store.Set(ctx, "lastrun:"+taskID, t.Unix())
}
