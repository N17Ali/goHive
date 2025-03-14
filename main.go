package main

import (
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID       string
	Title    string
	TaskFunc func() error
	Interval time.Duration
	NextRun  time.Time
	CreateAt time.Time
	UpdateAt time.Time
}

func (t *Task) runTask() {
	for {
		time.Sleep(time.Until(t.NextRun))

		err := t.TaskFunc()
		if err != nil {
			fmt.Printf("[%s] Task failed %v\n", t.ID, err)
		} else {
			fmt.Printf("[%s] Task completed succesfully\n", t.ID)
		}

		t.NextRun = time.Now().Add(t.Interval)
	}
}

func main() {
	Task := &Task{
		ID:    "t1",
		Title: "Simple Task",
		TaskFunc: func() error {
			if time.Now().Second()%2 == 0 {
				return errors.New("task failed")
			}
			return nil
		},
		Interval: 10 * time.Second,
		NextRun:  time.Now().Add(10 * time.Second),
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	go Task.runTask()

	select {}
}
