package actions

import (
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/gin-gonic/gin"
)

func DeselectAllDo(ctx *gin.Context) {
	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}

	project.DeselectAll()

	state.SaveProject(ctx, con, project)

	ctx.JSON(http.StatusOK, gin.H{})
}

func DeselectAllUndo(ctx *gin.Context) {
	type Params struct {
		SelectedWires    []types.IDType `json:"selectedWireIDs"`
		SelectedCircuits []types.IDType `json:"selectedCircuitIDs"`
	}
	var params Params

	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	var currentScene = project.GetCurrentScene()

	for _, id := range params.SelectedWires {
		if !currentScene.HasWire(id) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "no wire with specified id",
				"id":    id,
			})
			return
		}
	}

	for _, id := range params.SelectedCircuits {
		if !currentScene.HasCircuit(id) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "no circuit with specified id",
				"id":    id,
			})
			return
		}
	}

	for _, id := range params.SelectedWires {
		project.SelectedWires[id] = true
	}
	for _, id := range params.SelectedCircuits {
		project.SelectedCircuits[id] = true
	}

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})
}
