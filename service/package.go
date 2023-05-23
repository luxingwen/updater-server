package service


import (
	"github.com/google/uuid"
	"updater-server/model"
	"updater-server/pkg/app"
)

type PackageService struct {}

func (ps *PackageService) GetAllPackages(ctx *app.Context, versionUUID string) ([]model.Package, error) {
	var packages []model.Package
	err := ctx.DB.Model(&model.Package{}).Where("version_uuid = ?", versionUUID).Find(&packages).Error
	
	if err != nil {
		return nil, err
	}

	return packages, nil
}


func (ps *PackageService) CreatePackage(ctx *app.Context, mpackager *model.Package) error {
	mpackager.Uuid = uuid.New().String()
	result := ctx.DB.Create(&mpackager)
	return result.Error
}

func (ps *PackageService) UpdatePackage(ctx *app.Context, id string, updatedPackage *model.Package) error {
	return nil
}

func (ps *PackageService) DeletePackage(ctx *app.Context, id string) error {
	result := ctx.DB.Delete(&model.Package{}, "id = ?", id)
	return result.Error
}
