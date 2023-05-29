package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"

	"github.com/gin-gonic/gin"

	"encoding/json"

	"github.com/google/uuid"
)

type ProgramActionController struct {
	Service                    *service.ProgramActionService
	ClientService              *service.ClientService
	TaskService                *service.TaskService
	TaskExecutionRecordService *service.TaskExecutionRecordService
}

func (pac *ProgramActionController) GetProgramActionByUUID(c *app.Context) {

	var query model.ProgramAction

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	action, err := pac.Service.GetProgramActionByUUID(c, query.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(action)
}

func (pac *ProgramActionController) GetAllProgramActions(c *app.Context) {
	var query model.ReqProgramActionQuery

	err := c.ShouldBindJSON(&query)
	if err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())

		return
	}

	actions, err := pac.Service.GetAllProgramActions(c, query.ProgramUuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(actions)
}

func (pac *ProgramActionController) CreateProgramAction(c *app.Context) {
	var action model.ProgramAction
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pac.Service.CreateProgramAction(c, &action)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(action)
}

func (pac *ProgramActionController) UpdateProgramAction(c *app.Context) {
	var updatedAction model.ProgramAction
	uuid := c.Param("uuid")
	if err := c.ShouldBindJSON(&updatedAction); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pac.Service.UpdateProgramAction(c, uuid, &updatedAction)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(updatedAction)
}

func (pac *ProgramActionController) DeleteProgramAction(c *app.Context) {
	uuid := c.Param("uuid")
	err := pac.Service.DeleteProgramAction(c, uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Program Action deleted successfully"})
}

func (pac *ProgramActionController) CreateActionTask(c *app.Context) {
	var actionTask *model.ReqTaskProgramAction
	if err := c.ShouldBindJSON(&actionTask); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	clients, err := pac.ClientService.GetClientByHostInfo(c, actionTask.HostInfo)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	if len(clients) == 0 {
		c.JSONError(http.StatusInternalServerError, "no client found")
		return
	}

	action, err := pac.Service.GetProgramActionByUUID(c, actionTask.Content.ProgramActionUuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	content, err := json.Marshal(actionTask.Content)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	task := &model.Task{
		TaskID:      uuid.New().String(),
		TaskName:    actionTask.Name,
		Description: actionTask.Description,
		Creater:     actionTask.Creater,
		TaskType:    string(action.ActionType),
		Content:     string(content),
	}

	err = pac.TaskService.CreateTask(c, task)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	if actionTask.BatchTask.IsBatchTask() {
		taskBatchesInfoList := actionTask.BatchTask.GenerateTaskBatchesInfo(clients)

		for _, taskBatchInfo := range taskBatchesInfoList {
			bsContent, err := json.Marshal(taskBatchInfo)
			if err != nil {
				c.JSONError(http.StatusInternalServerError, err.Error())
				return
			}

			batchesTask := &model.Task{
				TaskID:     taskBatchInfo.TaskID,
				TaskType:   "batches",
				NextTaskID: taskBatchInfo.NextTaskID,
				Content:    string(bsContent),
			}

			err = pac.TaskService.CreateTask(c, batchesTask)
			if err != nil {
				c.JSONError(http.StatusInternalServerError, err.Error())
				return
			}

			for _, client := range taskBatchInfo.Clients {
				err = pac.createTaskExecutionRecord(c, taskBatchInfo.TaskID, client, action.Content, action.ActionType)
				if err != nil {
					c.JSONError(http.StatusInternalServerError, err.Error())
					return
				}

				if action.ActionType == model.ActionTypeComposite {
					err = pac.createSubTaskExecutionRecords(c, task, client, taskBatchInfo.TaskID, action)
					if err != nil {
						c.JSONError(http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}
	} else {
		for _, client := range clients {
			err = pac.createTaskExecutionRecord(c, task.TaskID, client, action.Content, action.ActionType)
			if err != nil {
				c.JSONError(http.StatusInternalServerError, err.Error())
				return
			}

			if action.ActionType == model.ActionTypeComposite {
				err = pac.createSubTaskExecutionRecords(c, task, client, "", action)
				if err != nil {
					c.JSONError(http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
	}

	c.JSONSuccess(task)
}

func (pac *ProgramActionController) createTaskExecutionRecord(c *app.Context, taskID string, client string, content string, taskType model.ActionType) error {
	taskExecutionRecord := &model.TaskExecutionRecord{
		RecordID:   uuid.New().String(),
		TaskID:     taskID,
		ClientUUID: client,
		Content:    content,
		TaskType:   string(taskType),
	}

	err := pac.TaskExecutionRecordService.CreateRecord(c, taskExecutionRecord)
	if err != nil {
		return err
	}

	return nil
}

func (pac *ProgramActionController) createSubTaskExecutionRecords(c *app.Context, task *model.Task, client string, parentRecordID string, action *model.ProgramAction) error {
	actionTemplates := make([]*model.TemplateAction, 0)
	err := json.Unmarshal([]byte(action.Content), &actionTemplates)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return err
	}

	for _, actionTemplate := range actionTemplates {
		subAction, err := pac.Service.GetProgramActionByUUID(c, actionTemplate.ProgramActionUuid)
		if err != nil {
			return err
		}

		subTaskExecutionRecord := &model.TaskExecutionRecord{
			RecordID:       actionTemplate.Uuid,
			TaskID:         task.TaskID,
			ClientUUID:     client,
			Content:        subAction.Content,
			ParentRecordID: parentRecordID,
		}

		err = pac.TaskExecutionRecordService.CreateRecord(c, subTaskExecutionRecord)
		if err != nil {
			return err
		}
	}

	return nil
}
