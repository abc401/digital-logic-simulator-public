package actions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/math"
	"github.com/gin-gonic/gin"
)

type CreateCircuitParams struct {
	CircuitID   types.IDType ``
	CircuitType string       `binding:"required"`
	LocScr      math.Vec2    `binding:"required"`
}

func CreateCircuitDo(ctx *gin.Context) {
	var params CreateCircuitParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	currentScene := project.GetCurrentScene()

	newCircuit, ok := state.DefaultCircuits[params.CircuitType]

	if !ok {
		var sceneID = -1
		var icType = ""
		for id, ic := range project.ICs {
			if !strings.EqualFold(ic, params.CircuitType) {
				continue
			}
			sceneID = int(id)
			icType = ic
		}

		if sceneID == -1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":        "Invalid circuit type.",
				"circuit-type": params.CircuitType,
			})
		}

		var targetScene = project.Scenes[types.IDType(sceneID)]

		var icInputsID, _ = targetScene.ICInputs.Unwrap()
		var ICInputs = targetScene.GetCircuit(icInputsID)

		var icOutputsID, _ = targetScene.ICOutputs.Unwrap()
		var ICOutputs = targetScene.GetCircuit(icOutputsID)

		newCircuit = types.Circuit{
			CircuitType: icType,
			NInputPins:  ICInputs.NInputPins,
			NOutputPins: ICOutputs.NOutputPins,
			Props:       types.CircuitProps{},
		}
	}

	newCircuit.ID = params.CircuitID
	newCircuit.PosWrl = project.View.ScreenToWorld(params.LocScr)
	newCircuit.InputWires = make([]*types.Wire, newCircuit.NInputPins)
	newCircuit.OutputWires = make([]*types.Wire, newCircuit.NOutputPins)

	if err := currentScene.AddCircuit(params.CircuitID, &newCircuit); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
			"id":    params.CircuitID,
		})
		return
	}
	state.SaveProject(ctx, con, project)
	helpers.PrintCurrentScene(project)
	fmt.Printf("Created Circuit: %s\n", helpers.SPrettyPrint(newCircuit))
	ctx.JSON(http.StatusOK, gin.H{})
}

func CreateCircuitUndo(ctx *gin.Context) {
	var params CreateCircuitParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	currentScene := project.GetCurrentScene()

	helpers.PrintCurrentScene(project)

	if err := currentScene.DeleteCircuit(params.CircuitID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":      "Circuit does not exist.",
			"circuit-id": params.CircuitID,
		})
		return
	}

	state.SaveProject(ctx, con, project)
	helpers.PrintCurrentScene(project)
	ctx.JSON(http.StatusOK, gin.H{})
}
