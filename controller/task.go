package controller

import (
	"encoding/json"
	"net/http"
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
