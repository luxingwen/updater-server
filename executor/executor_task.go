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

	// 如果状态是暂停、停止或者是运行中
	if taskInfo.TaskStatus == "paused" || taskInfo.TaskStatus == "stopped" || taskInfo.TaskStatus == "running" {
		return nil
	}

	// 判断状态是否完成
	if taskInfo.TaskStatus == "completed" || recordInfo.Status == "failed" || recordInfo.Status == "success" {

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

	getTaskIdAndType := func(tInfo model.TaskContentInfo) (string, string) {
		if tInfo.TaskID != "" {
			return tInfo.TaskID, "task"
		}
		return tInfo.TaskRecordId, "record"
	}

	tcontent := taskContent.Content.([]model.TaskContentInfo)

	// 任务执行
	if taskInfo.TaskType == "batches" {

		for _, tInfo := range tcontent {

			taskID, taskType := getTaskIdAndType(tInfo)
			taskExecItem := TaskExecItem{
				TaskID:   taskID,
				Category: taskContent.Type,
				TaskType: taskType,
			}

			err = EnqueueTask(ctx, taskExecItem)
			if err != nil {
				ctx.Logger.Error("enqueue task error:", err)
				return err
			}
		}
	} else {

		taskContentInfo := tcontent[0]

		// 任务执行
		taskID, taskType := getTaskIdAndType(taskContentInfo)

		taskExecItem := TaskExecItem{
			TaskID:   taskID,
			Category: taskContent.Type,
			TaskType: taskType,
		}

		err = EnqueueTask(ctx, taskExecItem)
		if err != nil {
			ctx.Logger.Error("enqueue task error:", err)
			return err
		}
	}

	// 更新任务状态

	err = es.TaskService.UpdateTaskStatus(es.WsContext.AppContext(), task.TaskID, "running")

	return err
}
