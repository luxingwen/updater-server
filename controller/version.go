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

// @summary Get all versions
// @description Get all versions
// @tags version
// @accept json
// @produce json
// @param version query string false "name of the version to get"
// @param programId query string false "programId of the version to get"
// @router /v1/version/list [post]
func (vc *VersionController) GetAllVersions(c *app.Context) {
	// Implementation goes here...

	var query model.ReqVersionQuery
	if err := c.ShouldBindQuery(&query); err != nil {
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

// @summary Get version by UUID
// @description Get version by UUID
// @tags version
// @accept json
// @produce json
// @param uuid path string true "UUID of the version to get"
// @router /v1/version/get/{uuid} [post]
func (vc *VersionController) GetVersionByUUID(c *app.Context) {
	// Implementation goes here...
}

// @summary Create version
// @description Create new version
// @tags version
// @accept json
// @produce json
// @param body body model.Version true "version object that needs to be added"
// @router /v1/version/create [post]
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

// @summary Update version
// @description Update existing version
// @tags version
// @accept json
// @produce json
// @param uuid path string true "UUID of the version to update"
// @param body body model.Version true "version object that needs to be updated"
// @router /v1/version/update/{uuid} [post]
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

// @summary Delete version
// @description Delete version by UUID
// @tags version
// @accept json
// @produce json
// @param uuid path string true "UUID of the version to delete"
// @router /v1/version/delete/{uuid} [post]
func (vc *VersionController) DeleteVersion(c *app.Context) {
	// Implementation goes here...
	uuid := c.Param("uuid")
	err := vc.Service.DeleteVersion(c, uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Version deleted successfully"})
}
