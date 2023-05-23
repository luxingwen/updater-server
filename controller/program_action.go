package controller

import (
	"net/http"
	"updater-server/model"
	"updater-server/service"
	"updater-server/pkg/app"

	"github.com/gin-gonic/gin"
)

type ProgramActionController struct {
	Service *service.ProgramActionService
}

func (pac *ProgramActionController) GetAllProgramActions(c *app.Context) {
	programUUID := c.Param("programUUID")
	actions, err := pac.Service.GetAllProgramActions(c, programUUID)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(actions)
}

func (pac *ProgramActionController) CreateProgramAction(c *app.Context) {
	var action model.ProgramAction
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pac.Service.CreateProgramAction(c, &action)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(action)
}

func (pac *ProgramActionController) UpdateProgramAction(c *app.Context) {
	var updatedAction model.ProgramAction
	uuid := c.Param("uuid")
	if err := c.ShouldBindJSON(&updatedAction); err != nil {
		c.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := pac.Service.UpdateProgramAction(c, uuid, &updatedAction)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSONSuccess(updatedAction)
}

func (pac *ProgramActionController) DeleteProgramAction(c *app.Context) {
	uuid := c.Param("uuid")
	err := pac.Service.DeleteProgramAction(c, uuid)
	if err != nil {
		c.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSONSuccess(gin.H{"message": "Program Action deleted successfully"})
}
