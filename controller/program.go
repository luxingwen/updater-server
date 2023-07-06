package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProgramController struct {
	Service *service.ProgramService
}

// 获取所有程序信息
// @Tags program
// @Summary 获取所有程序信息
// @Description 获取所有程序信息
// @Accept json
// @Produce json
// @Param query body model.ReqProgrameQuery true "查询参数"
// @Success 200 {object} app.Response "Success"
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

// 删除程序
// @Tags program
// @Summary 删除程序
// @Description 删除程序
// @Accept json
// @Produce json
// @Param query body model.ReqProgrameQuery true "查询参数"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/program/delete [post]
func (pc *ProgramController) DeleteProgram(c *app.Context) {

	var query model.ReqProgrameQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pc.Service.DeleteProgram(c, query.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Program deleted successfully"})
}

// 创建程序
// @Tags program
// @Summary 创建程序
// @Description 创建程序
// @Accept json
// @Produce json
// @Param body body model.Program true "program object that needs to be created"
// @Success 200 {object} app.Response "Success"
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

// 更新程序
// @Tags program
// @Summary 更新程序
// @Description 更新程序
// @Accept json
// @Produce json
// @Param body body model.Program true "program object that needs to be updated"
// @Success 200 {object} app.Response "Success"
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

// 获取程序详情信息
// @Tags program
// @Summary 获取程序详情信息
// @Description 获取程序详情信息
// @Accept json
// @Produce json
// @Param query body model.ReqProgrameQuery true "查询参数"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/program/detail [post]
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
