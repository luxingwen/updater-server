package controller

import (
	"fmt"
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

	reqbody, err := json.Marshal(actionTask)
	if err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	clients, err := pac.ClientService.GetClientByHostInfo(c, actionTask.HostInfo)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	if len(clients) == 0 {
		c.Logger.Errorf("no client found, HostInfo: %v", actionTask.HostInfo)
		c.JSONError(http.StatusInternalServerError, "no client found")
		return
	}

	action, err := pac.Service.GetProgramActionByUUID(c, actionTask.Content.ProgramActionUuid)
	if err != nil {
		c.Logger.Errorf("GetProgramActionByUUID error:%s, ProgramActionUuid: %s", err.Error(), actionTask.Content.ProgramActionUuid)
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	task := &model.Task{
		TaskID:      uuid.New().String(),
		TaskName:    actionTask.Name,
		Description: actionTask.Description,
		Creater:     actionTask.Creater,
		TaskType:    action.Name,
		TaskStatus:  model.TaskStatusPreparing,
		Ext:         string(reqbody),
	}

	err = pac.TaskService.CreateTask(c, task)
	if err != nil {
		c.Logger.Errorf("CreateTask error:%s, Task: %v", err.Error(), task)
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	taskContentInfoList := make([]model.TaskContentInfo, 0)

	taskContent := &model.TaskContent{
		Type: "task",
	}

	if actionTask.BatchTask.IsBatchTask() {

		taskBatchesInfoList := actionTask.BatchTask.GenerateTaskBatchesInfo(clients)

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

			err = pac.TaskService.CreateTask(c, batchesTask)
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

				err = pac.createTaskExecutionRecord(c, taskBatchInfo.TaskID, recordId, client, action)
				if err != nil {
					c.Logger.Errorf("createTaskExecutionRecord error:%s, TaskBatchInfo: %v", err.Error(), taskBatchInfo)
					c.JSONError(http.StatusInternalServerError, err.Error())
					return
				}

			}
			batchesTaskContent := &model.TaskContent{
				Type:    "record",
				Content: batchesTaskContentInfoList,
			}
			err = pac.TaskService.UpdateTaskContent(c, batchesTask.TaskID, batchesTaskContent)
			if err != nil {
				c.Logger.Errorf("UpdateTaskContent error:%s, TaskBatchInfo: %v", err.Error(), taskBatchInfo)
				c.JSONError(http.StatusInternalServerError, err.Error())
				return
			}

		}
	} else {
		taskContent.Type = "record"
		for _, client := range clients {

			recordId := uuid.New().String()
			taskContentInfoList = append(taskContentInfoList, model.TaskContentInfo{
				TaskRecordId: recordId,
				Sequence:     0,
			})

			err = pac.createTaskExecutionRecord(c, task.TaskID, recordId, client, action)
			if err != nil {
				c.Logger.Errorf("createTaskExecutionRecord error:%s, Task: %v", err.Error(), task)
				c.JSONError(http.StatusInternalServerError, err.Error())
				return
			}

		}
	}

	taskContent.Content = taskContentInfoList

	err = pac.TaskService.UpdateTaskContent(c, task.TaskID, taskContent)
	if err != nil {
		c.Logger.Errorf("UpdateTaskContent error:%s, Task: %v", err.Error(), task)
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(task)
}

func (pac *ProgramActionController) createTaskExecutionRecord(c *app.Context, taskID string, recordId string, client string, action *model.ProgramAction) error {

	taskExecutionRecord := &model.TaskExecutionRecord{
		RecordID:   recordId,
		TaskID:     taskID,
		ClientUUID: client,
		Content:    pac.getTaskRecordContent(c, action),
		TaskType:   string(action.ActionType),
		Name:       action.Name,
		Category:   "sub",
		Status:     model.TaskStatusPreparing,
	}

	if action.ActionType == model.ActionTypeComposite {
		taskExecutionRecord.Category = "root"
	}

	err := pac.TaskExecutionRecordService.CreateRecord(c, taskExecutionRecord)
	if err != nil {
		return err
	}

	if action.ActionType == model.ActionTypeComposite {
		err = pac.createSubTaskExecutionRecords(c, taskID, client, recordId, action)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pac *ProgramActionController) createSubTaskExecutionRecords(c *app.Context, taskId string, client string, parentRecordID string, action *model.ProgramAction) error {
	actionTemplates := make([]*model.TemplateAction, 0)
	err := json.Unmarshal([]byte(action.Content), &actionTemplates)
	if err != nil {
		c.Logger.Errorf("json.Unmarshal error:%s, Content: %s", err.Error(), action.Content)
		c.JSONError(http.StatusInternalServerError, err.Error())
		return err
	}

	model.GenerateTemplateActionNextUuids(actionTemplates)

	contentList := make([]model.TaskContentInfo, 0)

	for i, actionTemplate := range actionTemplates {

		recordId := actionTemplate.TaskRecordId

		subAction, err := pac.Service.GetProgramActionByUUID(c, actionTemplate.Uuid)
		if err != nil {
			return err
		}
		subTaskExecutionRecord := &model.TaskExecutionRecord{
			RecordID:       recordId,
			TaskID:         taskId,
			ClientUUID:     client,
			Content:        pac.getTaskRecordContent(c, subAction),
			ParentRecordID: parentRecordID,
			Category:       "sub",
			TaskType:       string(action.ActionType),
			Name:           action.Name,
			NextRecordID:   actionTemplate.NextTaskRecordId,
			Status:         model.TaskStatusPreparing,
		}

		contentList = append(contentList, model.TaskContentInfo{
			TaskRecordId: recordId,
			Sequence:     i + 1,
		})

		if subAction.ActionType == model.ActionTypeComposite {
			subTaskExecutionRecord.Category = "root"
		}

		err = pac.TaskExecutionRecordService.CreateRecord(c, subTaskExecutionRecord)
		if err != nil {
			return err
		}

		if subAction.ActionType == model.ActionTypeComposite {
			err = pac.createSubTaskExecutionRecords(c, taskId, client, recordId, subAction)
			if err != nil {
				return err
			}
		}
	}

	taskContent := &model.TaskContent{
		Type:    "record",
		Content: contentList,
	}

	pac.TaskExecutionRecordService.UpdaterRecordContent(c, parentRecordID, taskContent)
	return nil
}

func (pac *ProgramActionController) getTaskRecordContent(c *app.Context, action *model.ProgramAction) (r string) {
	taskContent := model.TaskContent{
		Type:    "program_script",
		Content: action.Content,
	}

	if action.ActionType == model.ActionTypeDownload {
		taskContent.Type = "program_download"
	}

	if action.ActionType == model.ActionTypeComposite {
		taskContent.Type = "record"
	}

	b, _ := json.Marshal(taskContent)
	return string(b)
}
