package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @title Program API
// @version 1.0
// @description This is the Program API

// ProgramController ...
type ProgramController struct {
	Service *service.ProgramService
}

// @Summary Get all programs
// @Description Get all programs
// @Tags program
// @Accept json
// @Produce json
// @Param programName query string false "name of the program to get"
// @Param teamId query string false "teamId of the program to get"
// @Router /api/v1/program/list [post]
func (pc *ProgramController) GetAllPrograms(c *app.Context) {
	var query model.ReqProgrameQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	response, err := pc.Service.GetAllPrograms(c, &query)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(response)
}

// @Summary Delete program
// @Description Delete program by UUID
// @Tags program
// @Accept json
// @Produce json
// @Param uuid path string true "UUID of the program to delete"
// @Router /api/v1/program/delete [post]
func (pc *ProgramController) DeleteProgram(c *app.Context) {
	uuid := c.Param("uuid")
	err := pc.Service.DeleteProgram(c, uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Program deleted successfully"})
}

// @Summary Create program
// @Description Create new program
// @Tags program
// @Accept json
// @Produce json
// @Param body body model.Program true "program object that needs to be added"
// @Router /api/v1/program/create [post]
func (pc *ProgramController) CreateProgram(c *app.Context) {
	var program model.Program
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	program.Uuid = uuid.New().String()
	err := pc.Service.CreateProgram(c, &program)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(program)
}

// @Summary Update program
// @Description Update existing program
// @Tags program
// @Accept json
// @Produce json
// @Param body body model.Program true "program object that needs to be updated"
// @Router /api/v1/program/update [post]
func (pc *ProgramController) UpdateProgram(c *app.Context) {
	var updatedProgram model.Program
	if err := c.ShouldBindJSON(&updatedProgram); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.Service.UpdateProgram(c, &updatedProgram)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(updatedProgram)
}

func (pc *ProgramController) GetProgramDetail(c *app.Context) {
	var query model.ReqProgrameQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	response, err := pc.Service.GetProgram(c, query.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(response)
}
