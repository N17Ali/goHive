package task

import "time"

type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Interval    time.Duration `json:"inteval"`
}
