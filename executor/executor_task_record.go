package executor

import (
	"encoding/json"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

// 执行任务记录
func (es *ExecutorServer) ExecuteTaskRecord(ctx *app.Context, task TaskExecItem) (err error) {
	recordInfo, err := service.NewTaskExecutionRecordService().GetRecordInfo(ctx, task.TaskID)
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
