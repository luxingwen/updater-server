package service

import (
	"fmt"
	"time"
	"updater-server/model"
	"updater-server/pkg/app"
)

type FileInfoService struct {
}

func NewFileInfoService() *FileInfoService {
	return &FileInfoService{}
}

// 创建文件信息
func (s *FileInfoService) CreateFileInfo(ctx *app.Context, fileInfo *model.FileInfo) error {
	fileInfo.CreateAt = time.Now()
	fileInfo.UpdateAt = time.Now()

	err := ctx.DB.Create(fileInfo).Error
	if err != nil {
		return fmt.Errorf("failed to create fileInfo: %w", err)
	}
	return nil
}

// 更新文件信息
func (s *FileInfoService) UpdateFileInfo(ctx *app.Context, fileInfo *model.FileInfo) error {
	fileInfo.UpdateAt = time.Now()

	err := ctx.DB.Model(&model.FileInfo{}).Where("uuid = ?", fileInfo.Uuid).Updates(fileInfo).Error
	if err != nil {
		return fmt.Errorf("failed to update fileInfo: %w", err)
	}
	return nil
}

// 获取文件信息
func (s *FileInfoService) GetFileInfo(ctx *app.Context, uuid string) (*model.FileInfo, error) {
	var fileInfo model.FileInfo
	err := ctx.DB.Model(&model.FileInfo{}).Where("uuid = ?", uuid).First(&fileInfo).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get fileInfo: %w", err)
	}

	return &fileInfo, nil
}

// 删除文件信息
func (s *FileInfoService) DeleteFileInfo(ctx *app.Context, uuid string) error {

	err := ctx.DB.Model(&model.FileInfo{}).Where("uuid = ?", uuid).Delete(&model.FileInfo{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete fileInfo: %w", err)
	}

	return nil
}

// 查询文件信息
func (s *FileInfoService) QueryFileInfo(ctx *app.Context, param model.ReqFileQuery) (r *model.PagedResponse, err error) {

	var fileInfo []model.FileInfo
	var count int64

	// 查询条件
	db := ctx.DB.Model(&model.FileInfo{}).Where("1 = 1")
	if param.Name != "" {
		db = db.Where("name like ?", "%"+param.Name+"%")
	}

	if param.Type != "" {
		db = db.Where("type = ?", param.Type)
	}

	if param.Creater != "" {
		db = db.Where("creater = ?", param.Creater)
	}

	if param.TeamId != "" {
		db = db.Where("team_id = ?", param.TeamId)
	}

	if param.DirUuid != "" {
		db = db.Where("parent_id = ?", param.DirUuid)
	}

	// 查询
	err = db.Offset(param.GetOffset()).Limit(param.PageSize).Find(&fileInfo).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query fileInfo: %w", err)
	}

	// 查询总数
	err = db.Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count fileInfo: %w", err)
	}

	return &model.PagedResponse{
		Total:    count,
		Current:  param.Current,
		PageSize: param.PageSize,
		Data:     fileInfo,
	}, nil

}
