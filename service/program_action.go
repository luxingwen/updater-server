package service

import (
	"updater-server/model"
	"updater-server/pkg/app"

	"github.com/google/uuid"
)

type ProgramActionService struct{}

func (pas *ProgramActionService) GetAllProgramActions(ctx *app.Context, programUUID string) ([]model.ProgramAction, error) {
	var actions []model.ProgramAction
	err := ctx.DB.Model(&model.ProgramAction{}).Where("program_uuid = ?", programUUID).Find(&actions).Error
	if err != nil {
		return nil, err
	}
	return actions, nil
}

func (pas *ProgramActionService) CreateProgramAction(ctx *app.Context, action *model.ProgramAction) error {
	action.Uuid = uuid.New().String()
	result := ctx.DB.Create(&action)
	return result.Error
}

func (pas *ProgramActionService) UpdateProgramAction(ctx *app.Context, uuid string, updatedAction *model.ProgramAction) error {
	var action model.ProgramAction
	result := ctx.DB.First(&action, "uuid = ?", uuid)

	if result.Error != nil {
		return result.Error
	}

	result = ctx.DB.Save(&updatedAction)
	return result.Error
}

func (pas *ProgramActionService) DeleteProgramAction(ctx *app.Context, uuid string) error {
	result := ctx.DB.Delete(&model.ProgramAction{}, "uuid = ?", uuid)
	return result.Error
}

func (pas *ProgramActionService) GetProgramActionByUUID(ctx *app.Context, uuid string) (*model.ProgramAction, error) {
	var action model.ProgramAction
	result := ctx.DB.First(&action, "uuid = ?", uuid)
	return &action, result.Error
}
