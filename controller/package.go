package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PackageController is the controller for managing packages.
type PackageController struct {
	Service *service.PackageService
}

// GetAllPackages retrieves all packages.
// @Tags Packages
// @Produce json
// @Param Pagination query model.Pagination true "Pagination data"
// @Param package query string false "Package Name"
// @Param arch query string false "Architecture"
// @Router /v1/packages/list [post]
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

// CreatePackage creates a new package.
// @Summary Create a new package
// @Tags Packages
// @Produce json
// @Param Package body model.Package true "Package data"
// @Router /v1/packages/create [post]
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

// UpdatePackage updates a package by ID.
// @Tags Packages
// @Produce json
// @Param id path string true "Package ID"
// @Param Package body model.Package true "Package data"
// @Router /v1/packages/update/{id} [post]
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

// DeletePackage deletes a package by ID.
// @Tags Packages
// @Produce json
// @Param id path string true "Package ID"
// @Router /v1/packages/delete/{id} [post]
func (pc *PackageController) DeletePackage(c *app.Context) {
	id := c.Param("id")
	err := pc.Service.DeletePackage(c, id)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Package deleted successfully"})
}
