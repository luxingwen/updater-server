package service

import (
	"fmt"
	"time"
	"updater-server/model"
	"updater-server/pkg/app"
)

type PreTaskService struct {
}

func NewPreTaskService() *PreTaskService {
	return &PreTaskService{}
}

func (s *PreTaskService) CreatePreTask(ctx *app.Context, preTask *model.PreTask) error {
	preTask.CreateAt = time.Now()
	preTask.UpdateAt = time.Now()

	err := ctx.DB.Create(preTask).Error
	if err != nil {
		return fmt.Errorf("failed to create preTask: %w", err)
	}
	return nil
}

// 更新预设任务
func (s *PreTaskService) UpdatePreTask(ctx *app.Context, preTask *model.PreTask) error {
	preTask.UpdateAt = time.Now()

	err := ctx.DB.Save(preTask).Error
	if err != nil {
		return fmt.Errorf("failed to update preTask: %w", err)
	}
	return nil
}

// 删除预设任务
func (s *PreTaskService) DeletePreTask(ctx *app.Context, uuid string) error {
	err := ctx.DB.Delete(&model.PreTask{}, "uuid = ?", uuid).Error
	if err != nil {
		return fmt.Errorf("failed to delete preTask: %w", err)
	}
	return nil
}

// 获取预设任务
func (s *PreTaskService) GetPreTaskByUUID(ctx *app.Context, uuid string) (*model.PreTask, error) {
	var preTask model.PreTask
	err := ctx.DB.First(&preTask, "uuid = ?", uuid).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get preTask by uuid: %w", err)
	}
	return &preTask, nil
}

// 获取预设任务列表
func (s *PreTaskService) QueryPreTaskList(ctx *app.Context, param *model.ReqPresetTaskQuery) (r *model.PagedResponse, err error) {
	var preTasks []model.PreTask
	var total int64
	db := ctx.DB.Model(&model.PreTask{})
	if param.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", param.Name))
	}

	if param.Category != "" {
		db = db.Where("category = ?", param.Category)
	}

	if param.Type != "" {
		db = db.Where("type = ?", param.Type)
	}

	err = db.Offset(param.GetOffset()).Limit(param.PageSize).Find(&preTasks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query preTask list: %w", err)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query preTask list count: %w", err)
	}
	return &model.PagedResponse{
		Total:    total,
		Current:  param.Current,
		PageSize: param.PageSize,
		Data:     preTasks,
	}, nil
}
