package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"updater-server/executor"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"

	"github.com/google/uuid"
)

type TaskController struct {
	Service                    *service.TaskService
	TaskExecutionRecordService *service.TaskExecutionRecordService
}

// 查询任务
// @Tags task
// @Summary 查询任务
// @Description 查询任务
// @Accept json
// @Produce json
// @Param query body model.ReqTaskQuery true "查询参数"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/task/list [post]
func (tc *TaskController) QueryTasks(c *app.Context) {
	var query model.ReqTaskQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	r, err := tc.Service.GetAllTasks(c, &query)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(r)
}

// 获取任务信息
// @Tags task
// @Summary 获取任务信息
// @Description 获取任务信息
// @Accept json
// @Produce json
// @Param query body model.ReqTaskInfoParam true "查询参数"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/task/info [post]
func (tc *TaskController) GetTaskInfo(c *app.Context) {
	var query model.Task
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	r, err := tc.Service.GetTaskInfo(c, query.TaskID)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(r)
}

// 创建单个任务
// @Tags task
// @Summary 创建单个任务
// @Description 创建单个任务
// @Accept json
// @Produce json
// @Param task body model.ReqTaskSingleCreate true "任务信息"
// @Success 200 {string} app.Response "Success"
// @Router /api/v1/task/create/single [post]
func (tc *TaskController) CreateSingleTask(c *app.Context) {

	var task model.ReqTaskSingleCreate
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	taskContent := model.TaskContent{
		Type: task.Type,
	}

	if task.Type == "script" {
		taskContent.Content = task.Script
	}

	if task.Type == "file" {
		taskContent.Content = task.DownloadFile
	}

	taskContentBytes, err := json.Marshal(taskContent)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	record := model.TaskExecutionRecord{
		TaskID:     uuid.New().String(),
		RecordID:   uuid.New().String(),
		Content:    string(taskContentBytes),
		ClientUUID: task.ClientUuid,
		Status:     model.TaskStatusPreparing,
		CreatedAt:  time.Now(),
		Timeout:    task.GetTimeout(),
	}

	if err := tc.TaskExecutionRecordService.CreateRecord(c, &record); err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	executor.EnqueueTask(c, executor.TaskExecItem{
		TaskID:   record.RecordID,
		Category: "record",
		TaskType: "sub",
		TraceId:  c.TraceID,
	})

	mdata := make(map[string]interface{})
	mdata["record_id"] = record.RecordID

	c.JSONSuccess(mdata)

}

// 创建多个任务
// @Tags task
// @Summary 创建多个任务
// @Description 创建多个任务
// @Accept json
// @Produce json
// @Param task body model.ReqTaskMultiCreate true "任务信息"
// @Success 200 {string} app.Response "Success"
// @Router /api/v1/task/create/multiple [post]
func (tc *TaskController) CreateMultipleTask(c *app.Context) {

	var param model.ReqTaskMultiCreate
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	task := model.Task{
		TaskID:      uuid.New().String(),
		TaskName:    param.TaskName,
		Description: param.Description,
		Creater:     param.Creater,
		TaskType:    param.Type,
		TaskStatus:  model.TaskStatusPreparing,
	}

	err := tc.Service.CreateTask(c, &task)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	recordContent := model.TaskContent{
		Type: param.Type,
	}

	if param.Type == "script" {
		recordContent.Content = param.Script
	}

	if param.Type == "file" {
		recordContent.Content = param.DownloadFile
	}

	recordContentBytes, err := json.Marshal(recordContent)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	taskContent := &model.TaskContent{
		Type: "record",
	}

	taskContentInfoList := make([]model.TaskContentInfo, 0)

	for _, item := range param.ClientUuids {

		record := model.TaskExecutionRecord{
			TaskID:     task.TaskID,
			RecordID:   uuid.New().String(),
			Content:    string(recordContentBytes),
			ClientUUID: item,
			Status:     model.TaskStatusPreparing,
			CreatedAt:  time.Now(),
			Timeout:    param.GetContentTimeout(),
		}

		if err := tc.TaskExecutionRecordService.CreateRecord(c, &record); err != nil {
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		taskContentInfoList = append(taskContentInfoList, model.TaskContentInfo{
			TaskRecordId: record.RecordID,
			Sequence:     0,
		})
	}

	taskContent.Content = taskContentInfoList

	err = tc.Service.UpdateTaskContent(c, task.TaskID, taskContent)
	if err != nil {
		c.Logger.Errorf("UpdateTaskContent error:%s, Task: %v", err.Error(), task)
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	executor.EnqueueTask(c, executor.TaskExecItem{
		TaskID:   task.TaskID,
		Category: "task",
		TaskType: "root",
		TraceId:  c.TraceID,
	})

	c.JSONSuccess(task)
}

// 创建批次任务
// @Tags task
// @Summary 创建批次任务
// @Description 创建批次任务
// @Accept json
// @Produce json
// @Param task body model.ReqTaskBatchCreate true "任务信息"
// @Success 200 {string} app.Response "Success"
// @Router /api/v1/task/create/batch [post]
func (tc *TaskController) CreateBatchTask(c *app.Context) {

	var param model.ReqTaskBatchCreate
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	task := model.Task{
		TaskID:      uuid.New().String(),
		TaskName:    param.TaskName,
		Description: param.Description,
		Creater:     param.Creater,
		TaskType:    param.Type,
		TaskStatus:  model.TaskStatusPreparing,
		TraceId:     c.TraceID,
	}

	err := tc.Service.CreateTask(c, &task)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	recordContent := model.TaskContent{
		Type: param.Type,
	}

	if param.Type == "script" {
		recordContent.Content = param.Script
	}

	if param.Type == "file" {
		recordContent.Content = param.DownloadFile
	}

	recordContentBytes, err := json.Marshal(recordContent)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	taskContent := &model.TaskContent{
		Type: "record",
	}

	taskContentInfoList := make([]model.TaskContentInfo, 0)

	taskBatchesInfoList := param.BatchTask.GenerateTaskBatchesInfo(param.ClientUuids)

	for _, taskBatchInfo := range taskBatchesInfoList {
		bsContent, err := json.Marshal(taskBatchInfo)
		if err != nil {
			c.Logger.Errorf("json.Marshal error:%s, TaskBatchInfo: %v", err.Error(), taskBatchInfo)
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		batchTaskName := fmt.Sprintf("第%d批次", taskBatchInfo.Sequence)

		batchesTask := &model.Task{
			TaskID:       taskBatchInfo.TaskID,
			TaskType:     "batches",
			NextTaskID:   taskBatchInfo.NextTaskID,
			Content:      string(bsContent),
			ParentTaskID: task.TaskID,
			Category:     "root",
			TaskName:     batchTaskName,
			TaskStatus:   model.TaskStatusPreparing,
		}

		taskContentInfoList = append(taskContentInfoList, model.TaskContentInfo{
			TaskID:   taskBatchInfo.TaskID,
			Sequence: taskBatchInfo.Sequence,
		})

		err = tc.Service.CreateTask(c, batchesTask)
		if err != nil {
			c.Logger.Errorf("CreateTask error:%s, Task: %v", err.Error(), batchesTask)
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		batchesTaskContentInfoList := make([]model.TaskContentInfo, 0)

		for _, client := range taskBatchInfo.Clients {

			recordId := uuid.New().String()

			batchesTaskContentInfoList = append(batchesTaskContentInfoList, model.TaskContentInfo{
				TaskRecordId: recordId,
				Sequence:     0,
			})

			record := model.TaskExecutionRecord{
				TaskID:     batchesTask.TaskID,
				RecordID:   recordId,
				Content:    string(recordContentBytes),
				ClientUUID: client,
				Status:     model.TaskStatusPreparing,
				CreatedAt:  time.Now(),
				Timeout:    param.GetContentTimeout(),
			}

			if err := tc.TaskExecutionRecordService.CreateRecord(c, &record); err != nil {
				c.JSONError(http.StatusInternalServerError, err.Error())
				return
			}

		}
		batchesTaskContent := &model.TaskContent{
			Type:    "record",
			Content: batchesTaskContentInfoList,
		}
		err = tc.Service.UpdateTaskContent(c, batchesTask.TaskID, batchesTaskContent)
		if err != nil {
			c.Logger.Errorf("UpdateTaskContent error:%s, TaskBatchInfo: %v", err.Error(), taskBatchInfo)
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
	}
	taskContent.Content = taskContentInfoList

	err = tc.Service.UpdateTaskContent(c, task.TaskID, taskContent)
	if err != nil {
		c.Logger.Errorf("UpdateTaskContent error:%s, Task: %v", err.Error(), task)
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	executor.EnqueueTask(c, executor.TaskExecItem{
		TaskID:   task.TaskID,
		Category: "task",
		TaskType: "root",
		TraceId:  c.TraceID,
	})

	c.JSONSuccess(task)
}

func (tc *TaskController) CreateFileDownload(c *app.Context) {
	var taskFile model.FileDownloadTask
	if err := c.ShouldBindJSON(&taskFile); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	reqbody, err := json.Marshal(taskFile)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	task := &model.Task{
		TaskID:      uuid.New().String(),
		TaskName:    taskFile.Name,
		Description: "",
		Creater:     taskFile.Creater,
		TaskType:    "file_download",
		TaskStatus:  model.TaskStatusPreparing,
		Ext:         string(reqbody),
	}

	err = tc.Service.CreateTask(c, task)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	taskContent := &model.TaskContent{
		Type:    "file_download",
		Content: taskFile.Content,
	}
	taskContentStr, err := json.Marshal(taskContent)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	for _, item := range taskFile.Clients {

		taskExecutionRecord := &model.TaskExecutionRecord{
			RecordID:   uuid.New().String(),
			TaskID:     task.TaskID,
			ClientUUID: item,
			Content:    string(taskContentStr),
			TaskType:   "file_download",
			Name:       "",
			Category:   "sub",
			Status:     model.TaskStatusPreparing,
		}

		err = tc.TaskExecutionRecordService.CreateRecord(c, taskExecutionRecord)
		if err != nil {
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}
	}

}
