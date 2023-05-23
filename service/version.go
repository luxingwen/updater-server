package service


import (
	"github.com/google/uuid"
	"updater-server/model"
	"updater-server/pkg/app"
)

type VersionService struct {}

func (vs *VersionService) GetAllVersions(ctx *app.Context, query *model.ReqVersionQuery) (*model.PagedResponse, error) {
	// Implementation goes here...
	return nil, nil
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
