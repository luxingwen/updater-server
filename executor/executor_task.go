package executor

import (
	"encoding/json"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

// 执行任务
func (es *ExecutorServer) ExecuteTask(ctx *app.Context, task TaskExecItem) (err error) {
	taskInfo, err := service.NewTaskService().GetTaskInfo(ctx, task.TaskID)
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
