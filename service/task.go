package service

import (
	"encoding/json"
	"updater-server/model"
	"updater-server/pkg/app"

	"gorm.io/gorm"
)

type TaskService struct{}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (ts *TaskService) CreateTask(ctx *app.Context, task *model.Task) error {
	tx := ctx.DB.Begin() // 开始事务
	err := ctx.DB.Create(task).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	tx.Commit() // 提交事务
	return nil
}

func (ts *TaskService) UpdateTask(ctx *app.Context, updatedTask *model.Task) error {
	result := ctx.DB.Model(&model.Task{}).Where("task_id = ?", updatedTask.TaskID).Updates(updatedTask)
	return result.Error
}

// 更新任务状态
func (ts *TaskService) UpdateTaskStatus(ctx *app.Context, taskID string, status string) error {
	result := ctx.DB.Model(&model.Task{}).Where("task_id = ?", taskID).Update("task_status", status)
	return result.Error
}

func (ts *TaskService) UpdateTaskContent(ctx *app.Context, taskID string, content interface{}) error {

	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	result := ctx.DB.Model(&model.Task{}).Where("task_id = ?", taskID).Update("content", string(b))
	return result.Error
}

func (ts *TaskService) DeleteTask(ctx *app.Context, taskID string) error {
	var task model.Task
	result := ctx.DB.Where("task_id = ?", taskID).Delete(&task)
	return result.Error
}

func (ts *TaskService) GetTaskInfo(ctx *app.Context, taskID string) (*model.Task, error) {
	var task model.Task
	result := ctx.DB.Where("task_id = ?", taskID).First(&task)
	return &task, result.Error
}

func (ts *TaskService) GetAllTasks(ctx *app.Context, query *model.ReqTaskQuery) (*model.PagedResponse, error) {
	sess := ctx.DB.Session(&gorm.Session{})

	if query.TaskName != "" {
		sess = sess.Where("task_name = ?", query.TaskName)
	}

	if query.TeamId != "" {
		sess = sess.Where("team_id = ?", query.TeamId)
	}

	if len(query.TaskIds) > 0 {
		sess = sess.Where("task_id IN(?)", query.TaskIds)
	}

	var tasks []model.Task
	result := sess.Order("created_at DESC").Limit(query.PageSize).Offset(query.GetOffset()).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	ctx.DB.Model(&model.Task{}).Count(&total)

	response := &model.PagedResponse{
		Data:     tasks,
		Current:  query.Current,
		PageSize: query.PageSize,
		Total:    total,
	}

	return response, nil
}

// 检查任务状态，检查任务是否完成
func (ts *TaskService) CheckTaskStatus(ctx *app.Context, taskID string) (bool, error) {
	taskInfo, err := ts.GetTaskInfo(ctx, taskID)
	if err != nil {
		return false, err
	}

	if taskInfo.TaskStatus == "completed" || taskInfo.TaskStatus == "failed" || taskInfo.TaskStatus == "success" {
		return true, nil
	}

	taskContent := &model.TaskContent{}

	err = json.Unmarshal([]byte(taskInfo.Content), taskContent)
	if err != nil {
		ctx.Logger.Error("unmarshal task content error:", err)
		return false, err
	}

	tcontent := taskContent.Content.([]model.TaskContentInfo)
	if taskContent.Type == "task" {

		for _, item := range tcontent {

			isDone, err := ts.CheckTaskStatus(ctx, item.TaskID)
			if err != nil {
				return false, err
			}

			if !isDone {
				return false, nil
			}
		}

		// 跟新状态
		err = ts.UpdateTaskStatus(ctx, taskID, "completed")
		return true, nil
	}

	recordService := NewTaskExecutionRecordService()

	if taskContent.Type == "record" {

		for _, item := range tcontent {

			isDone, err := recordService.CheckTaskStatus(ctx, item.TaskRecordId)
			if err != nil {
				return false, err
			}

			if !isDone {
				return false, nil
			}

		}

		// 更新状态
		err = ts.UpdateTaskStatus(ctx, taskID, "completed")
		return true, nil
	}

	return false, nil

}
