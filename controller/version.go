package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"

	"github.com/gin-gonic/gin"
)

type VersionController struct {
	Service *service.VersionService
}

// 获取所有版本信息
// @Tags version
// @Summary 获取所有版本信息
// @Description 获取所有版本信息
// @Accept json
// @Produce json
// @Param query body model.ReqVersionQuery true "查询参数"
// @Success 200 {object} app.Response "Success"
// @Router /v1/version/list [post]
func (vc *VersionController) GetAllVersions(c *app.Context) {
	// Implementation goes here...

	var query model.ReqVersionQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	response, err := vc.Service.GetAllVersions(c, &query)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(response)
}

// 获取版本详细信息
// @Tags version
// @Summary 获取版本详细信息
// @Description 获取版本详细信息
// @Accept json
// @Produce json
// @Param query body model.Version true "查询参数"
// @Success 200 {object} app.Response "Success"
// @Router /v1/version/info [post]
func (vc *VersionController) GetVersionByUUID(c *app.Context) {
	// Implementation goes here...

	var query model.Version

	err := c.ShouldBindJSON(&query)
	if err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	response, err := vc.Service.GetVersionInfo(c, query.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(response)
}

// 创建版本
// @Tags version
// @Summary 创建版本
// @Description 创建版本
// @Accept json
// @Produce json
// @Param body body model.Version true "version object that needs to be added"
// @Success 200 {object} app.Response "Success"
// @Router /v1/version/create [post]
func (vc *VersionController) CreateVersion(c *app.Context) {
	// Implementation goes here...

	var version model.Version
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := vc.Service.CreateVersion(c, &version)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(version)
}

// 更新版本
// @Tags version
// @Summary 更新版本
// @Description 更新版本
// @Accept json
// @Produce json
// @Param body body model.Version true "version object that needs to be updated"
// @Success 200 {object} app.Response "Success"
// @Router /v1/version/update [post]
func (vc *VersionController) UpdateVersion(c *app.Context) {
	// Implementation goes here...

	var updatedVersion model.Version
	if err := c.ShouldBindJSON(&updatedVersion); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := vc.Service.UpdateVersion(c, &updatedVersion)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(updatedVersion)
}

// 删除版本
// @Tags version
// @Summary 删除版本
// @Description 删除版本
// @Accept json
// @Produce json
// @Param body body model.Version true "version object that needs to be deleted"
// @Success 200 {object} app.Response "Success"
// @Router /v1/version/delete [post]
func (vc *VersionController) DeleteVersion(c *app.Context) {
	var query model.Version

	err := c.ShouldBindJSON(&query)
	if err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	err = vc.Service.DeleteVersion(c, query.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Version deleted successfully"})
}
