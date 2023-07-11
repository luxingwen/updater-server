package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

type DangerousCommandController struct {
	DangerousCommandService *service.DangerousCommandService
}

// 创建危险指令
// @Summary 创建危险指令
// @Description 创建危险指令
// @Tags 危险指令
// @Accept json
// @Produce json
// @Param param body model.DangerousCommand true "危险指令参数"
// @Success 200 {object} model.DangerousCommandInfoResponse
// @Router /api/v1/dangerous_command/create [post]
func (dc *DangerousCommandController) CreateDangerousCommand(ctx *app.Context) {
	var param = &model.DangerousCommand{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	err := dc.DangerousCommandService.CreateDangerousCommand(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// 更新危险指令
// @Summary 更新危险指令
// @Description 更新危险指令
// @Tags 危险指令
// @Accept json
// @Produce json
// @Param param body model.DangerousCommand true "危险指令参数"
// @Success 200 {object} model.DangerousCommandInfoResponse
// @Router /api/v1/dangerous_command/update [post]
func (dc *DangerousCommandController) UpdateDangerousCommand(ctx *app.Context) {
	var param = &model.DangerousCommand{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	err := dc.DangerousCommandService.UpdateDangerousCommand(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// 删除危险指令
// @Summary 删除危险指令
// @Description 删除危险指令
// @Tags 危险指令
// @Accept json
// @Produce json
// @Param param body model.ReqUuidParam true "危险指令参数"
// @Success 200 {object} app.Response "Success"
// @Router /api/v1/dangerous_command/delete [post]
func (dc *DangerousCommandController) DeleteDangerousCommand(ctx *app.Context) {
	var param = &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	err := dc.DangerousCommandService.DeleteDangerousCommand(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(nil)
}

// 获取危险指令列表
// @Summary 获取危险指令列表
// @Description 获取危险指令列表
// @Tags 危险指令
// @Accept json
// @Produce json
// @Param param body model.ReqDangerousCommandQuery true "危险指令参数"
// @Success 200 {object} model.DangerousCommandQueryResponse
// @Router /api/v1/dangerous_command/list [post]
func (dc *DangerousCommandController) GetDangerousCommandList(ctx *app.Context) {
	var param = &model.ReqDangerousCommandQuery{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	r, err := dc.DangerousCommandService.GetDangerousCommandList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}

// 获取危险指令详情
// @Summary 获取危险指令详情
// @Description 获取危险指令详情
// @Tags 危险指令
// @Accept json
// @Produce json
// @Param param body  model.ReqUuidParam  true "危险指令uuid"
// @Success 200 {object}  model.DangerousCommandInfoResponse
// @Router /api/v1/dangerous_command/detail [post]
func (dc *DangerousCommandController) GetDangerousCommandInfo(ctx *app.Context) {
	var param = &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	r, err := dc.DangerousCommandService.GetDangerousCommandInfo(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}
