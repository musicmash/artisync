package tasks

import "errors"

var (
	ErrInvalidUUID  = errors.New("task_id should be valid uuid")
	ErrTaskNotFound = errors.New("task_id not found")
)
