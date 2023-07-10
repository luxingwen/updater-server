package controller

import (
	"encoding/json"
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

type PreTaskController struct {
	PreTaskService *service.PreTaskService
}

// @Summary 创建预设任务
// @Description 创建预设任务
// @Tags pre_task
// @Accept json
// @Produce json
// @Param param body model.ReqPresetTaskCreate true "创建参数"
// @Success 200 {object} model.PreTaskInfoResponse
// @Router /api/v1/pre_task/create [post]
func (ptc *PreTaskController) CreatePreTask(c *app.Context) {
	var param model.ReqPresetTaskCreate
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	if param.Category != "mutilTask" && param.Category != "batchesTask" && param.Category != "programActionTask" {
		c.JSONError(http.StatusBadRequest, "invalid category")
		return
	}

	preTask := model.PreTask{
		Category: param.Category,
	}

	param.SetPreTask(&preTask)

	if preTask.Content == "" {
		c.JSONError(http.StatusBadRequest, "invalid content")
		return
	}

	if err := ptc.PreTaskService.CreatePreTask(c, &preTask); err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(preTask)
}

// @Summary 查询预设任务
// @Description 查询预设任务
// @Tags pre_task
// @Accept json
// @Produce json
// @Param query body model.ReqPresetTaskQuery true "查询参数"
// @Success 200 {object} model.PreTaskQueryResponse
// @Router /api/v1/pre_task/list [post]
func (ptc *PreTaskController) QueryPreTasks(c *app.Context) {
	var query model.ReqPresetTaskQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	r, err := ptc.PreTaskService.QueryPreTaskList(c, &query)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(r)
}

// @Summary 删除预设任务
// @Description 删除预设任务
// @Tags pre_task
// @Accept json
// @Produce json
// @Param query body model.ReqUuidParam true "预设任务ID"
// @Success 200 {object} app.Response
// @Router /api/v1/pre_task/delete [post]
func (ptc *PreTaskController) DeletePreTask(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	if err := ptc.PreTaskService.DeletePreTask(c, param.Uuid); err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(nil)
}

// @Summary 更新预设任务
// @Description 更新预设任务
// @Tags pre_task
// @Accept json
// @Produce json
// @Param param body model.ReqPresetTaskUpdate true "预设任务信息"
// @Success 200 {object} model.PreTaskInfoResponse
// @Router /api/v1/pre_task/update [post]
func (ptc *PreTaskController) UpdatePreTask(c *app.Context) {
	var param model.ReqPresetTaskUpdate
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	if param.Category != "mutilTask" && param.Category != "batchesTask" && param.Category != "programActionTask" {
		c.JSONError(http.StatusBadRequest, "invalid category")
		return
	}

	preTask, err := ptc.PreTaskService.GetPreTaskByUUID(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	param.SetPreTask(preTask)

	if preTask.Content == "" {
		c.JSONError(http.StatusBadRequest, "invalid content")
		return
	}

	if err := ptc.PreTaskService.UpdatePreTask(c, preTask); err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(preTask)
}

// @Summary 获取预设任务详细信息
// @Description 获取预设任务详细信息
// @Tags pre_task
// @Accept json
// @Produce json
// @Param query body model.ReqUuidParam true "预设任务ID"
// @Success 200 {object} model.PreTaskInfoResponse
// @Router /api/v1/pre_task/detail [post]
func (ptc *PreTaskController) GetPreTaskDetail(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	preTask, err := ptc.PreTaskService.GetPreTaskByUUID(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(preTask)
}

// 执行预设任务
// @Summary 执行预设任务
// @Description 执行预设任务
// @Tags pre_task
// @Accept json
// @Produce json
// @Param param body model.ReqUuidParam true "预设任务ID"
// @Success 200 {object} model.TaskInfoResponse
// @Router /api/v1/pre_task/execute [post]
func (ptc *PreTaskController) ExecutePreTask(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	preTask, err := ptc.PreTaskService.GetPreTaskByUUID(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	if preTask.Category == "mutilTask" {
		taskController := TaskController{
			Service:                    &service.TaskService{},
			TaskExecutionRecordService: &service.TaskExecutionRecordService{},
		}

		var mutilTaskParam model.ReqTaskMultiCreate

		err = json.Unmarshal([]byte(preTask.Content), &mutilTaskParam)
		if err != nil {
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		taskController.createMultipleTask(c, mutilTaskParam)
		return
	}

	if preTask.Category == "batchesTask" {
		taskController := TaskController{
			Service:                    &service.TaskService{},
			TaskExecutionRecordService: &service.TaskExecutionRecordService{},
		}

		var batchesTaskParam model.ReqTaskBatchCreate

		err = json.Unmarshal([]byte(preTask.Content), &batchesTaskParam)
		if err != nil {
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		taskController.createBatchTask(c, batchesTaskParam)
		return
	}

	c.JSONError(http.StatusInternalServerError, "没有实现执行器的预设任务类型")
}
