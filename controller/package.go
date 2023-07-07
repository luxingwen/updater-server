package controller

import (
	"net/http"
	"os"
	"path"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PackageController struct {
	Service *service.PackageService
}

// 获取程序所有的安装包
// @Tags Packages
// @Summary 获取程序所有的安装包
// @Description 获取程序所有的安装包
// @Accept json
// @Produce json
// @Param query body model.ReqPackageQuery true "查询参数"
// @Success 200 {object} model.GetAllProgramPackageResponse
// @Router /api/v1/program/package/list [post]
func (pc *PackageController) GetAllPackages(c *app.Context) {
	var query model.ReqPackageQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	response, err := pc.Service.GetAllPackages(c, query.VersionUuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(response)
}

// 创建安装包
// @Tags Packages
// @Summary 创建安装包
// @Description 创建安装包
// @Accept json
// @Produce json
// @Param Package body model.Package true "Package data"
// @Success 200 {object} model.CreateProgramPackageResponse
// @Router /api/v1/program/package/create [post]
func (pc *PackageController) CreatePackage(c *app.Context) {
	var mpackage model.Package
	if err := c.ShouldBindJSON(&mpackage); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	mpackage.Uuid = uuid.New().String()
	err := pc.Service.CreatePackage(c, &mpackage)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(mpackage)
}

// 更新安装包
// @Tags Packages
// @Summary 更新安装包
// @Description 更新安装包
// @Accept json
// @Produce json
// @Param Package body model.Package true "Package data"
// @Success 200 {object} model.CreateProgramPackageResponse
// @Router /api/v1/program/package/update/{id} [post]
func (pc *PackageController) UpdatePackage(c *app.Context) {
	var updatedPackage model.Package
	id := c.Param("id")
	if err := c.ShouldBindJSON(&updatedPackage); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.Service.UpdatePackage(c, id, &updatedPackage)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(updatedPackage)
}

// 删除安装包
// @Tags Packages
// @Summary 删除安装包
// @Description 删除安装包
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/program/package/delete/{id} [post]
func (pc *PackageController) DeletePackage(c *app.Context) {
	id := c.Param("id")
	err := pc.Service.DeletePackage(c, id)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Package deleted successfully"})
}

func (pc *PackageController) UploadFile(c *app.Context) {
	programUuid := c.Param("programUuid")

	savePathdir := path.Join(c.Config.PkgFileDir, programUuid)

	_, err := os.Stat(savePathdir)
	if err != nil && os.IsNotExist(err) {

		err := os.MkdirAll(savePathdir, 0755)
		if err != nil {
			c.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

	}

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	saveFilePath := path.Join(savePathdir, file.Filename)

	// 将文件保存到指定路径
	err = c.SaveUploadedFile(file, saveFilePath)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(saveFilePath)

}

// 删除文件
// @Tags Packages
// @Summary 删除文件
// @Description 删除文件
// @Accept json
// @Produce json
// @Param query body model.ReqDeletePackageFile true "File Name"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/program/package/file/delete/{programUuid} [post]
func (pc *PackageController) DeleteFile(c *app.Context) {
	programUuid := c.Param("programUuid")

	var query model.ReqDeletePackageFile

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	filePath := path.Join(c.Config.PkgFileDir, programUuid, query.FileName)

	err := os.Remove(filePath)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(gin.H{"message": "File deleted successfully"})
}
