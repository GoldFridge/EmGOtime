package models

import (
	"time"
)

// Duration is a custom type to represent time.Duration in a format Swagger can understand
type Duration string

// Task represents a task in the system
// @Description Task represents a task in the system
type Task struct {
	// ID is the unique identifier of the task
	ID uint `json:"id" gorm:"primary_key"`

	// UserID is the ID of the user to whom the task is assigned
	UserID uint `json:"user_id"`

	// Name is the name of the task
	Name string `json:"name"`

	// StartTime is the start time of the task
	StartTime time.Time `json:"start_time"`

	// EndTime is the end time of the task
	EndTime time.Time `json:"end_time"`

	// Duration is the duration of the task
	Duration Duration `json:"duration"`
}

func (t *Task) StartTask() {
	t.StartTime = time.Now()
}

// EndTask sets the end time of the task and calculates the duration
func (t *Task) EndTask() {
	t.EndTime = time.Now()
	duration := t.EndTime.Sub(t.StartTime)
	t.Duration = Duration(duration.String())
}
