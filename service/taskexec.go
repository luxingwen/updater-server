package service

import (
	"context"
	"encoding/json"
	"time"

	"updater-server/model"
	"updater-server/pkg/app"

	"github.com/go-redis/redis/v8"
)

const (
	TaskQueueKey = "task_queue"
)

func EnqueueTask(ctx *app.Context, task *model.Task) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return ctx.Redis.Enqueue(context.Background(), TaskQueueKey, string(taskJSON))
}

func ExecuteTask(ctx *app.Context) error {
	taskJSON, err := ctx.Redis.Dequeue(context.Background(), TaskQueueKey)
	if err != nil {
		return err
	}

	task := &model.Task{}
	err = json.Unmarshal([]byte(taskJSON), task)
	if err != nil {
		return err
	}

	if task.TaskType == "program" {

		
	}

	// 分批任务
	if task.TaskType == "batch" {

	}

	// Perform the task
	// You can define your own task execution logic here.
	// ...

	// If the task execution is successful, update the task status
	task.TaskStatus = "completed"
	err = ctx.DB.Save(&task).Error
	if err != nil {
		return err
	}

	return nil
}

func Worker(ctx context.Context, appCtx *app.Context) {

	defer func() {
		if r := recover(); r != nil {
			appCtx.Logger.Errorf("Recovered from panic in Worker: %v", r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// The context has been cancelled, stop the worker
			return
		default:
			// Continue to the next task
		}

		err := ExecuteTask(appCtx)
		if err != nil {
			if err == redis.Nil {
				// No task in the queue, sleep for a while and try again
				time.Sleep(time.Second)
				continue
			}
			// Log the error and continue to the next task
			appCtx.Logger.Errorf("Failed to execute task: %v", err)
			continue
		}

	}
}
