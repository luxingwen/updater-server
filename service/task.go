package service

import (
	"updater-server/model"
	"updater-server/pkg/app"

	"gorm.io/gorm"
)

type TaskService struct{}

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

func (ts *TaskService) DeleteTask(ctx *app.Context, taskID string) error {
	var task model.Task
	result := ctx.DB.Where("task_id = ?", taskID).Delete(&task)
	return result.Error
}

func (ts *TaskService) GetAllTasks(ctx *app.Context, query *model.ReqTaskQuery) (*model.PagedResponse, error) {
	sess := ctx.DB.Session(&gorm.Session{})

	if query.TaskName != "" {
		sess = sess.Where("task_name = ?", query.TaskName)
	}

	if query.TeamId != "" {
		sess = sess.Where("team_id = ?", query.TeamId)
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
