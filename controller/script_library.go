package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

type ScriptLibraryController struct {
	ScriptLibraryService *service.ScriptLibraryService
}

// 创建脚本库
// @Tags scriptLibrary
// @Summary 创建脚本库
// @Description 创建脚本库
// @Accept json
// @Produce json
// @Param param body model.ScriptLibrary true "脚本库"
// @Success 200 {object} model.ScriptLibraryInfoResponse
// @Router /api/v1/script_library/create [post]
func (slc *ScriptLibraryController) CreateScriptLibrary(c *app.Context) {
	var param model.ScriptLibrary
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := slc.ScriptLibraryService.CreateScriptLibrary(c, &param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(param)
}

// 更新脚本库
// @Tags scriptLibrary
// @Summary 更新脚本库
// @Description 更新脚本库
// @Accept json
// @Produce json
// @Param param body model.ScriptLibrary true "脚本库"
// @Success 200 {object} model.ScriptLibraryInfoResponse
// @Router /api/v1/script_library/update [post]
func (slc *ScriptLibraryController) UpdateScriptLibrary(c *app.Context) {
	var param model.ScriptLibrary
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := slc.ScriptLibraryService.UpdateScriptLibrary(c, &param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(param)
}

// 删除脚本库
// @Tags scriptLibrary
// @Summary 删除脚本库
// @Description 删除脚本库
// @Accept json
// @Produce json
// @Param param body model.ReqUuidParam true "脚本库ID"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/script_library/delete [post]
func (slc *ScriptLibraryController) DeleteScriptLibrary(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := slc.ScriptLibraryService.DeleteScriptLibrary(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess("ok")
}

// 获取脚本库详情
// @Tags scriptLibrary
// @Summary 获取脚本库详情
// @Description 获取脚本库详情
// @Accept json
// @Produce json
// @Param param body model.ReqUuidParam true "脚本库ID"
// @Success 200 {object} model.ScriptLibraryInfoResponse
// @Router /api/v1/script_library/detail [post]
func (slc *ScriptLibraryController) GetScriptLibrary(c *app.Context) {
	var param model.ReqUuidParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	response, err := slc.ScriptLibraryService.GetScriptLibrary(c, param.Uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(response)
}

// 获取脚本库列表
// @Tags scriptLibrary
// @Summary 获取脚本库列表
// @Description 获取脚本库列表
// @Accept json
// @Produce json
// @Param param body model.ReqScriptLibQuery true "脚本库列表"
// @Success 200 {object} model.ScriptLibraryQueryResponse
// @Router /api/v1/script_library/list [post]
func (slc *ScriptLibraryController) GetScriptLibraryList(c *app.Context) {
	var param model.ReqScriptLibQuery
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	response, err := slc.ScriptLibraryService.QueryScriptLibrary(c, param)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(response)
}
