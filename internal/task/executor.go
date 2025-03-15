package task

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/n17ali/gohive/internal/logs"
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
	err := runTaskFunction(task.ID)
	if err != nil {
		log.Printf("task %s failed: %v\n", task.Title, err)
		e.logger.LogTaskExecution(ctx, task.ID, "FAILED", err.Error())
	} else {
		log.Printf("task %s completed successfuly\n", task.Title)
		e.logger.LogTaskExecution(ctx, task.ID, "SUCCESS", "task executed successfuly")
	}
}

func runTaskFunction(taskID string) error {
	if time.Now().Second()%2 == 0 {
		return fmt.Errorf("simulated failure")
	}

	return nil
}
