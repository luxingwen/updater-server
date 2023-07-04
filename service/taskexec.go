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

func EnqueueTask(ctx *app.Context, task TaskExecItem) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return ctx.Redis.Enqueue(context.Background(), TaskQueueKey, string(taskJSON))
}

type TaskExecItem struct {
	TaskID   string `json:"task_id"`   // 任务ID
	Category string `json:"category"`  // 任务分类 task/record
	TaskType string `json:"task_type"` // 任务类型 root/batches
}

func Execute(ctx *app.Context) error {
	taskJSON, err := ctx.Redis.Dequeue(context.Background(), TaskQueueKey)
	if err != nil {
		return err
	}

	task := &TaskExecItem{}
	err = json.Unmarshal([]byte(taskJSON), task)
	if err != nil {
		return err
	}

	if task.Category == "task" {

	}

	// 分批任务
	if task.Category == "record" {

		recordInfo, err := NewTaskExecutionRecordService().GetRecordInfo(ctx, task.TaskID)
		if err != nil {
			ctx.Logger.Error("get record info error:", err)
			return err
		}

		if recordInfo.Status == "completed" {
			return nil
		}

	}
	return nil
}

// 执行任务记录
func ExecuteTaskRecord(ctx *app.Context, task TaskExecItem) (err error) {
	recordInfo, err := NewTaskExecutionRecordService().GetRecordInfo(ctx, task.TaskID)
	if err != nil {
		ctx.Logger.Error("get task info error:", err)
		return err
	}

	if recordInfo.Status == "completed" {
		if recordInfo.NextRecordID != "" {
			// 下一个任务
			nextTaskExecItem := TaskExecItem{
				TaskID:   recordInfo.NextRecordID,
				Category: task.Category,
				TaskType: task.TaskType,
			}

			err = EnqueueTask(ctx, nextTaskExecItem)
			if err != nil {
				ctx.Logger.Error("enqueue task error:", err)
				return err
			}
			return nil
		}
	}

	taskContent := &model.TaskContent{}

	err = json.Unmarshal([]byte(recordInfo.Content), taskContent)
	if err != nil {
		ctx.Logger.Error("unmarshal task content error:", err)
		return err
	}

	if taskContent.Type == "record" {
		tcontent := taskContent.Content.([]model.TaskContentInfo)
		taskExecItem := TaskExecItem{
			TaskID:   tcontent[0].TaskRecordId,
			Category: task.Category,
			TaskType: "",
		}

		err = EnqueueTask(ctx, taskExecItem)
		if err != nil {
			ctx.Logger.Error("enqueue task error:", err)
			return err
		}
	}

	if taskContent.Type == "script" {

	}

	return
}

func ExecuteTask(ctx *app.Context, task TaskExecItem) (err error) {
	taskInfo, err := NewTaskService().GetTaskInfo(ctx, task.TaskID)
	if err != nil {
		ctx.Logger.Error("get task info error:", err)
		return err
	}

	// 判断状态
	if taskInfo.TaskStatus == "completed" {

		if taskInfo.NextTaskID != "" {
			// 下一个任务
			nextTaskExecItem := TaskExecItem{
				TaskID:   taskInfo.NextTaskID,
				Category: task.Category,
				TaskType: task.TaskType,
			}
			err = EnqueueTask(ctx, nextTaskExecItem)
			if err != nil {
				ctx.Logger.Error("enqueue task error:", err)
				return err
			}
			return nil
		}

		//

		return nil
	}

	taskContent := &model.TaskContent{}

	err = json.Unmarshal([]byte(taskInfo.Content), taskContent)
	if err != nil {
		ctx.Logger.Error("unmarshal task content error:", err)
		return err
	}

	getTaskId := func(tInfo model.TaskContentInfo) string {
		if tInfo.TaskID != "" {
			return tInfo.TaskID
		}
		return tInfo.TaskRecordId
	}

	tcontent := taskContent.Content.([]model.TaskContentInfo)

	// 任务执行
	if taskInfo.TaskType == "batches" {

		for _, tInfo := range tcontent {

			taskExecItem := TaskExecItem{
				TaskID:   getTaskId(tInfo),
				Category: taskContent.Type,
				TaskType: "",
			}

			err = EnqueueTask(ctx, taskExecItem)
			if err != nil {
				ctx.Logger.Error("enqueue task error:", err)
				return err
			}
		}
		return
	}

	taskContentInfo := tcontent[0]

	// 任务执行

	taskExecItem := TaskExecItem{
		TaskID:   getTaskId(taskContentInfo),
		Category: taskContent.Type,
		TaskType: "",
	}

	err = EnqueueTask(ctx, taskExecItem)
	if err != nil {
		ctx.Logger.Error("enqueue task error:", err)
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

		err := Execute(appCtx)
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
