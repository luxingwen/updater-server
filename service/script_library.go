package service

import (
	"fmt"
	"time"
	"updater-server/model"
	"updater-server/pkg/app"
)

type ScriptLibraryService struct {
}

func NewScriptLibraryService() *ScriptLibraryService {
	return &ScriptLibraryService{}
}

func (s *ScriptLibraryService) CreateScriptLibrary(ctx *app.Context, scriptLibrary *model.ScriptLibrary) error {
	scriptLibrary.CreatAt = time.Now()
	scriptLibrary.UpdateAt = time.Now()

	err := ctx.DB.Create(scriptLibrary).Error
	if err != nil {
		return fmt.Errorf("failed to create scriptLibrary: %w", err)
	}
	return nil
}

// 更新脚本库
func (s *ScriptLibraryService) UpdateScriptLibrary(ctx *app.Context, scriptLibrary *model.ScriptLibrary) error {
	scriptLibrary.UpdateAt = time.Now()

	err := ctx.DB.Model(&model.ScriptLibrary{}).Where("uuid = ?", scriptLibrary.Uuid).Updates(scriptLibrary).Error
	if err != nil {
		return fmt.Errorf("failed to update scriptLibrary: %w", err)
	}
	return nil
}

// 获取脚本库
func (s *ScriptLibraryService) GetScriptLibrary(ctx *app.Context, uuid string) (*model.ScriptLibrary, error) {
	var scriptLibrary model.ScriptLibrary
	err := ctx.DB.Model(&model.ScriptLibrary{}).Where("uuid = ?", uuid).First(&scriptLibrary).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get scriptLibrary: %w", err)
	}

	return &scriptLibrary, nil
}

// 删除脚本库
func (s *ScriptLibraryService) DeleteScriptLibrary(ctx *app.Context, uuid string) error {
	err := ctx.DB.Model(&model.ScriptLibrary{}).Where("uuid = ?", uuid).Delete(&model.ScriptLibrary{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete scriptLibrary: %w", err)
	}

	return nil
}

// 查询脚本库
func (s *ScriptLibraryService) QueryScriptLibrary(ctx *app.Context, param model.ReqScriptLibQuery) (r *model.PagedResponse, err error) {

	var scriptLibraries []model.ScriptLibrary
	var total int64

	db := ctx.DB.Model(&model.ScriptLibrary{})

	if param.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", param.Name))
	}

	if param.Type != "" {
		db = db.Where("type = ?", param.Type)
	}

	if param.TeamId != "" {
		db = db.Where("team_id = ?", param.TeamId)
	}

	if param.Creater != "" {
		db = db.Where("creater = ?", param.Creater)
	}

	if param.Status > 0 {
		db = db.Where("status = ?", param.Status)
	}

	err = db.Offset(param.GetOffset()).Limit(param.PageSize).Find(&scriptLibraries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get scriptLibrary: %w", err)
	}

	err = db.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count scriptLibrary: %w", err)
	}

	return &model.PagedResponse{
		Total:    total,
		Current:  param.Current,
		PageSize: param.PageSize,
		Data:     scriptLibraries,
	}, nil

}
