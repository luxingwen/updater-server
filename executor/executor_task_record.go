package executor

import (
	"encoding/json"
	"strings"
	"time"
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

	// 如果状态是暂停、停止或者是运行中
	if recordInfo.Status == "paused" || recordInfo.Status == "stopped" || recordInfo.Status == "running" {
		ctx.Logger.Info("task status is:", recordInfo.Status)
		return nil
	}

	// 如果任务状态是已经完成
	if recordInfo.Status == "completed" || recordInfo.Status == "failed" || recordInfo.Status == "success" {
		if recordInfo.NextRecordID != "" {
			// 下一个任务
			nextTaskExecItem := TaskExecItem{
				TaskID:   recordInfo.NextRecordID,
				Category: task.Category,
				TaskType: task.TaskType,
				TraceId:  task.TraceId,
			}

			err = EnqueueTask(ctx, nextTaskExecItem)
			if err != nil {
				ctx.Logger.Error("enqueue task error:", err)
				return err
			}
			return nil
		}
		return
	}

	taskContent := &model.TaskContent{}

	err = json.Unmarshal([]byte(recordInfo.Content), taskContent)
	if err != nil {
		ctx.Logger.Error("unmarshal task content error:", err)
		es.ExecuteTaskRecordFailed(ctx, recordInfo, err.Error())
		return err
	}

	if taskContent.Type == "record" {
		// 更新状态
		err = es.TaskExecutionRecordService.UpdateRecordStatus(ctx, recordInfo.RecordID, "running")
		if err != nil {
			ctx.Logger.Error("update record status error:", err)

			return err
		}

		tcontent := taskContent.Content.([]model.TaskContentInfo)
		taskExecItem := TaskExecItem{
			TaskID:   tcontent[0].TaskRecordId,
			Category: task.Category,
			TaskType: "",
			TraceId:  task.TraceId,
		}

		err = EnqueueTask(ctx, taskExecItem)
		if err != nil {
			ctx.Logger.Error("enqueue task error:", err)
			es.ExecuteTaskRecordFailed(ctx, recordInfo, err.Error())
			return err
		}
	}

	ctx.Logger.Info("task content:", taskContent)
	if taskContent.Type == "program_script" {
		ctx.Logger.Info("script task")
		// 执行脚本任务
		es.ExecuteProgramScript(ctx, recordInfo)
		return
	}

	if taskContent.Type == "program_download" {
		ctx.Logger.Info("program_download task")
	}

	return err
}

// 执行程序下载
func (es *ExecutorServer) ExecuteProgramDownload(ctx *app.Context, recordInfo *model.TaskExecutionRecord) (err error) {

	return err
}

// 数据库里面存的信息
type TaskContentProgram struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type TaskContentProgramScript struct {
	Os      string `json:"os"`
	Content string `json:"content"`
}

type ScriptTaskRequest struct {
	TaskID      string   `json:"task_id"`
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	WorkDir     string   `json:"workDir"`
	Params      []string `json:"params"`
	Env         map[string]string
	Timeout     int    `json:"timeout"`
	Interpreter string `json:"interpreter"`
	Stdin       string `json:"stdin"`
}

// 执行程序脚本
func (es *ExecutorServer) ExecuteProgramScript(ctx *app.Context, recordInfo *model.TaskExecutionRecord) (err error) {

	clientInfo, err := es.ClientService.GetClientByUUID(ctx, recordInfo.ClientUUID)
	if err != nil {
		ctx.Logger.Error("get client error:", err)
		es.ExecuteTaskRecordFailed(ctx, recordInfo, err.Error())
		return err
	}

	tcontent := &TaskContentProgram{}
	err = json.Unmarshal([]byte(recordInfo.Content), tcontent)
	if err != nil {
		ctx.Logger.Error("unmarshal task content error:", err)
		es.ExecuteTaskRecordFailed(ctx, recordInfo, err.Error())
		return err
	}

	scriptContentList := make([]TaskContentProgramScript, 0)
	err = json.Unmarshal([]byte(tcontent.Content), &scriptContentList)
	if err != nil {
		ctx.Logger.Error("unmarshal task content error:", err)
		es.ExecuteTaskRecordFailed(ctx, recordInfo, err.Error())
		return err
	}

	if clientInfo.OS == "darwin" {
		clientInfo.OS = "linux"
	}

	scriptContent := ""
	for _, item := range scriptContentList {
		if strings.ToLower(item.Os) == strings.ToLower(clientInfo.OS) {
			ctx.Logger.Info("execute script:", item.Content)
			scriptContent = item.Content
			break
		}
	}

	if scriptContent == "" {
		ctx.Logger.Error("not found script content")
		es.ExecuteTaskRecordFailed(ctx, recordInfo, "not found script content")
		return
	}

	sq := ScriptTaskRequest{
		TaskID:  recordInfo.RecordID,
		Type:    "script",
		Content: scriptContent,
	}

	if sq.Timeout == 0 {
		sq.Timeout = 60
	}

	err = es.WsContext.SendRequest(recordInfo.ClientUUID, "v1/ExecuteScript", sq)

	if err != nil {
		ctx.Logger.Error("send request error:", err)
		es.ExecuteTaskRecordFailed(ctx, recordInfo, err.Error())
		return
	}

	mdata := make(map[string]interface{})
	mdata["status"] = "running"
	mdata["start_time"] = time.Now()
	err = es.TaskExecutionRecordService.UpdateRecordByMap(ctx, recordInfo.RecordID, mdata)
	if err != nil {
		ctx.Logger.Error("update record error:", err)
		return
	}

	return err
}

// 执行记录失败
func (es *ExecutorServer) ExecuteTaskRecordFailed(ctx *app.Context, recordInfo *model.TaskExecutionRecord, msg string) (err error) {
	mdata := make(map[string]interface{})
	mdata["status"] = "failed"
	mdata["message"] = msg
	mdata["start_time"] = time.Now()
	mdata["end_time"] = time.Now()
	err = es.TaskExecutionRecordService.UpdateRecordByMap(ctx, recordInfo.RecordID, mdata)
	if err != nil {
		ctx.Logger.Error("update record error:", err)
		return
	}
	return err
}
