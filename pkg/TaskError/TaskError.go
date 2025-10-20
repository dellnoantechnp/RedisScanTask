package TaskError

import "fmt"

// TaskError 任务错误结构体
type TaskError struct {
	Code int
	Err  error
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("TaskError: %d", e.Code)
}
