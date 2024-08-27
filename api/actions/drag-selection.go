package actions

import (
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/math"
	"github.com/gin-gonic/gin"
)

func DragSelectionDo(ctx *gin.Context) {
	type Params struct {
		DeltaWrl math.Vec2 `binding:"required"`
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

	for id := range project.SelectedCircuits {

		circuit := currentScene.GetCircuit(id)
		circuit.PosWrl = circuit.PosWrl.Add(params.DeltaWrl)
	}
	state.SaveProject(ctx, con, project)
	helpers.PrintCurrentScene(project)
	ctx.JSON(http.StatusOK, gin.H{})
}

func DragSelectionUndo(ctx *gin.Context) {
	type Params struct {
		DeltaWrl math.Vec2 `binding:"required"`
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

	for id := range project.SelectedCircuits {

		circuit := currentScene.GetCircuit(id)
		circuit.PosWrl = circuit.PosWrl.Sub(params.DeltaWrl)
	}
	state.SaveProject(ctx, con, project)
	helpers.PrintCurrentScene(project)
	ctx.JSON(http.StatusOK, gin.H{})
}
