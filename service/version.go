package service

import (
	"updater-server/model"
	"updater-server/pkg/app"

	"github.com/google/uuid"
)

type VersionService struct{}

func (vs *VersionService) GetAllVersions(ctx *app.Context, query *model.ReqVersionQuery) (*model.PagedResponse, error) {
	// Implementation goes here...

	var versions []model.Version
	result := ctx.DB.Model(&model.Version{}).Where("program_uuid = ?", query.ProgramUuid).Find(&versions)

	if result.Error != nil {
		return nil, result.Error
	}
	response := &model.PagedResponse{
		Data:     versions,
		Current:  query.Current,
		PageSize: query.PageSize,
		Total:    0,
	}
	return response, nil
}

func (vs *VersionService) CreateVersion(ctx *app.Context, version *model.Version) error {
	// Implementation goes here...
	version.Uuid = uuid.New().String()
	result := ctx.DB.Create(&version)
	return result.Error
}

func (vs *VersionService) UpdateVersion(ctx *app.Context, updatedVersion *model.Version) error {
	// Implementation goes here...

	result := ctx.DB.Model(&model.Version{}).Where("uuid = ?", updatedVersion.Uuid).Updates(updatedVersion)
	return result.Error
}

func (vs *VersionService) DeleteVersion(ctx *app.Context, uuid string) error {
	// Implementation goes here...

	result := ctx.DB.Delete(&model.Version{}, "uuid = ?", uuid)
	return result.Error
}

func (vs *VersionService) GetVersionInfo(ctx *app.Context, uuid string) (*model.Version, error) {
	var version model.Version
	result := ctx.DB.First(&version, "uuid = ?", uuid)
	return &version, result.Error
}
