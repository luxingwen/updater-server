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
