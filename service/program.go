package service

import (
	"updater-server/model"
	"updater-server/pkg/app"

	"gorm.io/gorm"
)

type ProgramService struct{}

func (ps *ProgramService) CreateProgram(ctx *app.Context, program *model.Program) error {
	tx := ctx.DB.Begin() // 开始事务

	// 创建程序
	result := tx.Create(program)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		return result.Error
	}

	// 创建初始化动作
	for _, action := range model.InitialActions {
		action.ProgramUUID = program.Uuid // 设置动作的 ProgramUUID
		result := tx.Create(&action)
		if result.Error != nil {
			tx.Rollback() // 回滚事务
			return result.Error
		}
	}

	tx.Commit() // 提交事务
	return nil
}

func (ps *ProgramService) UpdateProgram(ctx *app.Context, updatedProgram *model.Program) error {
	result := ctx.DB.Model(&model.Program{}).Where("uuid = ?", updatedProgram.Uuid).Updates(updatedProgram)
	return result.Error
}

func (ps *ProgramService) DeleteProgram(ctx *app.Context, uuid string) error {
	var program model.Program
	result := ctx.DB.Where("uuid = ?", uuid).Delete(&program)
	return result.Error
}

func (ps *ProgramService) GetAllPrograms(ctx *app.Context, query *model.ReqProgrameQuery) (*model.PagedResponse, error) {

	// Create a new session
	sess := ctx.DB.Session(&gorm.Session{})

	// Filter by program name if provided
	if query.ProgramName != "" {
		sess = sess.Where("name = ?", query.ProgramName)
	}

	// Filter by team ID if provided
	if query.TeamId != "" {
		sess = sess.Where("team_id = ?", query.TeamId)
	}

	var programs []model.Program
	result := sess.Limit(query.PageSize).Offset(query.GetOffset()).Find(&programs)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	ctx.DB.Model(&model.Program{}).Count(&total)

	response := &model.PagedResponse{
		Data:     programs,
		Current:  query.Current,
		PageSize: query.PageSize,
		Total:    total,
	}

	return response, nil
}

func (ps *ProgramService) GetProgram(ctx *app.Context, uuid string) (*model.Program, error) {
	var program model.Program
	result := ctx.DB.Where("uuid = ?", uuid).First(&program)
	if result.Error != nil {
		return nil, result.Error
	}
	return &program, nil
}
