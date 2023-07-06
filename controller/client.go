package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

type ClientController struct {
	Service *service.ClientService
}

// 查询客户端
// @Tags client
// @Summary 查询客户端
// @Description 查询客户端
// @Accept json
// @Produce json
// @Param query body model.ReqClientQuery true "查询参数"
// @Success 200 {object} app.Response "Success"
// @Router /v1/client/list [post]
func (cc *ClientController) QueryClients(c *app.Context) {
	var query model.ReqClientQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	r, err := cc.Service.QueryClient(c, &query)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(r)
}
