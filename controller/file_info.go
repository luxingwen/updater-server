package controller

import (
	"fmt"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

type FileInfoController struct {
	FileInfoService *service.FileInfoService
}

// 创建目录
// @Tags fileInfo
// @Summary 创建目录
// @Description 创建目录
// @Accept json
// @Produce json
// @Param param body model.ReqDirCreate true "目录路径"
// @Success 200 {object} model.FileInfoResponse
// @Router /api/v1/file_info/create_dir [post]
func (fic *FileInfoController) CreateDir(c *app.Context) {
	var param model.ReqDirCreate

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(400, err.Error())
		return
	}

	fileinfo := model.FileInfo{
		Path:    param.Dir,
		Creater: param.Creater,
		TeamId:  param.TeamId,
		IsDir:   true,
	}

	err := fic.FileInfoService.CreateFileInfo(c, &fileinfo)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess(fileinfo)
}

// 上传文件
// @Tags fileInfo
// @Summary 上传文件
// @Description 上传文件
// @Accept json
// @Produce json
// @Param file formData file true "文件"
// @Param dir formData string true "文件路径"
// @Param teamId formData string true "团队ID"
// @Param creater formData string true "创建者"
// @Success 200 {object} model.FileInfoResponse
// @Router /api/v1/file_info/upload [post]
func (fic *FileInfoController) UploadFile(c *app.Context) {
	var param model.ReqFileUpload
	if err := c.ShouldBind(&param); err != nil {
		c.JSONError(400, err.Error())
		return
	}

	// 保存文件

	saveFilename := fmt.Sprintf("%s/%s/%s", param.TeamId, param.Dir, param.File.Filename)

	err := c.SaveUploadedFile(param.File, saveFilename)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	fileinfo := model.FileInfo{
		Path:      param.Dir,
		TeamId:    param.TeamId,
		Creater:   param.Creater,
		StorePath: saveFilename,
		IsDir:     false,
	}

	err = fic.FileInfoService.CreateFileInfo(c, &fileinfo)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess(fileinfo)
}

// 删除文件
// @Tags fileInfo
// @Summary 删除文件
// @Description 删除文件
// @Accept json
// @Produce json
// @Param param body model.ReqUuidParam true "文件uuid"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/file_info/delete [post]
func (fic *FileInfoController) DeleteFile(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(400, err.Error())
		return
	}

	err := fic.FileInfoService.DeleteFileInfo(c, param.Uuid)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess("ok")
}

// 获取文件列表
// @Tags fileInfo
// @Summary 获取文件列表
// @Description 获取文件列表
// @Accept json
// @Produce json
// @Param param body model.ReqFileQuery true "目录路径"
// @Success 200 {object} model.FileInfoQueryResponse
// @Router /api/v1/file_info/list [post]
func (fic *FileInfoController) GetFileList(c *app.Context) {
	var param model.ReqFileQuery
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(400, err.Error())
		return
	}

	fileList, err := fic.FileInfoService.QueryFileInfo(c, param)
	if err != nil {
		c.JSONError(500, err.Error())
		return
	}
	c.JSONSuccess(fileList)
}
