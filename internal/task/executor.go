package task

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/n17ali/gohive/internal/logs"
)

const (
	MaxRetries   = 3
	RetryBackoff = 5 * time.Second
)

type TaskExecutor struct {
	interval time.Duration
	stopChan chan struct{}
	logger   logs.Logger
}

func NewTaskExecutor(interval time.Duration, logger logs.Logger) *TaskExecutor {
	return &TaskExecutor{
		interval: interval,
		stopChan: make(chan struct{}),
		logger:   logger,
	}
}

func (e *TaskExecutor) Start(ctx context.Context) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.runScheduledTasks(ctx)
		case <-e.stopChan:
			log.Println("task executor stopped.")
			return
		}
	}
}

func (e *TaskExecutor) runScheduledTasks(ctx context.Context) {
	log.Printf("geting all task...")
	tasks, err := GetAllTask(ctx)
	if err != nil {
		log.Printf("error fetching tasks:%v\n", err)
		return
	}
	now := time.Now()

	for _, task := range tasks {
		lastTaskRunTime, err := GetTaskLastRunTime(ctx, task.ID)
		if err != nil || now.Sub(lastTaskRunTime) >= task.Interval {
			go e.executeTask(ctx, task)
			SetTaskLatRunTime(ctx, task.ID, now)
		}
	}
}

func (e *TaskExecutor) executeTask(ctx context.Context, task Task) {
	var err error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		err = runTaskFunction(task.ID)
		if err == nil {
			log.Printf("task %s completed successfuly\n", task.Title)
			e.logger.LogTaskExecution(ctx, task.ID, "SUCCESS", "task executed successfuly")
			return
		}
		log.Printf("task %s failed on attempt %d: %v\n", task.Title, attempt, err)
		e.logger.LogTaskExecution(ctx, task.ID, "FAILED", "Retry attempt "+strconv.Itoa(attempt)+": "+err.Error())

		if attempt < MaxRetries {
			log.Printf("Retrying task %s in %v...\n", task.Title, RetryBackoff)
			time.Sleep(RetryBackoff)
		}
	}
	log.Printf("task %s failed after %d attempts\n", task.Title, MaxRetries)
	e.logger.LogTaskExecution(ctx, task.ID, "PERMANENT_FAILURE", "task failed after max retries")
}

func runTaskFunction(taskID string) error {
	if time.Now().Second()%2 == 0 {
		return fmt.Errorf("simulated failure")
	}

	return nil
}
