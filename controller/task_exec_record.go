package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

type TaskExecRecordController struct {
	Service *service.TaskExecutionRecordService
}

// 查询任务执行记录
// @Tags task_exec_record
// @Summary 查询任务执行记录
// @Description 查询任务执行记录
// @Accept json
// @Produce json
// @Param query body model.ReqTaskRecordeQuery true "查询参数"
// @Success 200 {object} model.TaskExecRecordQueryResponse
// @Router /api/v1/task/record/list [post]
func (tcr *TaskExecRecordController) QueryTaskExecRecords(c *app.Context) {
	var query model.ReqTaskRecordeQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	r, err := tcr.Service.GetAllTaskExecRecords(c, &query)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(r)
}

// 获取任务执行记录信息
// @Tags task_exec_record
// @Summary 获取任务执行记录信息
// @Description 获取任务执行记录信息
// @Accept json
// @Produce json
// @Param query body model.TaskRecordInfoParam true "查询参数"
// @Success 200 {object} model.TaskExecRecordInfoResponse
// @Router /api/v1/task/record/detail [post]
func (tcr *TaskExecRecordController) GetTaskExecRecordInfo(c *app.Context) {
	var query model.TaskRecordInfoParam
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	r, err := tcr.Service.GetRecordInfo(c, query.RecordId)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(r)
}
