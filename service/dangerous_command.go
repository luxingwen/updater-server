package service

import (
	"updater-server/model"

	"updater-server/pkg/app"
)

type DangerousCommandService struct {
}

// NewDangerousCommandService DangerousCommandService的构造函数
func NewDangerousCommandService() *DangerousCommandService {
	return &DangerousCommandService{}
}

// 创建危险指令
func (service *DangerousCommandService) CreateDangerousCommand(ctx *app.Context, param *model.DangerousCommand) error {
	return ctx.DB.Create(param).Error
}

// 更新危险指令
func (service *DangerousCommandService) UpdateDangerousCommand(ctx *app.Context, param *model.DangerousCommand) error {
	return ctx.DB.Model(param).Where("uuid = ?", param.Uuid).Updates(param).Error
}

// 删除危险指令
func (service *DangerousCommandService) DeleteDangerousCommand(ctx *app.Context, uuid string) error {
	return ctx.DB.Where("uuid = ?", uuid).Delete(&model.DangerousCommand{}).Error
}

// 获取危险指令列表
func (service *DangerousCommandService) GetDangerousCommandList(ctx *app.Context, param *model.ReqDangerousCommandQuery) (r *model.PagedResponse, err error) {

	var (
		data  []model.DangerousCommand
		total int64
	)

	db := ctx.DB.Model(&model.DangerousCommand{})

	if param.Name != "" {
		db = db.Where("name LIKE ?", "%"+param.Name+"%")
	}

	err = db.Offset(param.GetOffset()).Limit(param.PageSize).Find(&data).Error
	if err != nil {
		return
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	r = &model.PagedResponse{
		Current:  param.Current,
		PageSize: param.PageSize,
		Total:    total,
		Data:     data,
	}

	return
}

// 获取指令信息
func (service *DangerousCommandService) GetDangerousCommandInfo(ctx *app.Context, uuid string) (r *model.DangerousCommand, err error) {

	var data model.DangerousCommand

	err = ctx.DB.Where("uuid = ?", uuid).First(&data).Error
	if err != nil {
		return
	}

	r = &data

	return
}
