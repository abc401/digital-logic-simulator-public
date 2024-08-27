package actions

import (
	"fmt"
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/gin-gonic/gin"
)

type SwitchSceneParams struct {
	FromSceneID      types.IDType
	ToSceneID        types.IDType
	SelectedCircuits []types.IDType
	SelectedWires    []types.IDType
}

func SwitchSceneDo(ctx *gin.Context) {

	var params SwitchSceneParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	if _, ok := project.Scenes[params.ToSceneID]; !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no scene with specified id",
			"id":    params.ToSceneID,
		})
		return
	}
	project.CurrentSceneID = params.ToSceneID
	project.SelectedCircuits = map[types.IDType]bool{}

	project.ReEvaluateICs(params.ToSceneID, params.FromSceneID)

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})
	fmt.Printf("\n\nProject State: %s\n\n", helpers.SPrettyPrint(project))
}

func SwitchSceneUndo(ctx *gin.Context) {
	var params SwitchSceneParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	if _, ok := project.Scenes[params.FromSceneID]; !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no scene with specified id",
			"id":    params.ToSceneID,
		})
		return
	}

	project.CurrentSceneID = params.FromSceneID
	var currentScene = project.GetCurrentScene()

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
		if !currentScene.HasWire(id) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "no wire with specified id",
				"id":    id,
			})
			return
		}
	}

	for _, id := range params.SelectedCircuits {
		project.SelectedCircuits[id] = true
	}

	for _, id := range params.SelectedWires {
		project.SelectedWires[id] = true
	}

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})

	fmt.Printf("\n\nProject State: %s\n\n", helpers.SPrettyPrint(project))
}
