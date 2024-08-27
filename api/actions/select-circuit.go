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

type SelectCircuitParams struct {
	CircuitID types.IDType `json:"circuitID"`
}

func SelectCircuitDo(ctx *gin.Context) {
	var params SelectCircuitParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	var currentScene = project.GetCurrentScene()

	if !currentScene.HasObject(params.CircuitID) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No circuit with provided id",
			"id":    params.CircuitID,
		})
	}

	var circuit = currentScene.GetCircuit(params.CircuitID)

	for _, wire := range circuit.InputWires {
		if wire == nil {
			continue
		}
		if (wire.FromCircuit == circuit.ID && project.SelectedCircuits[wire.ToCircuit]) || (wire.ToCircuit == circuit.ID && project.SelectedCircuits[wire.FromCircuit]) {
			project.SelectedWires[wire.ID] = true
		}
	}
	for _, wire := range circuit.OutputWires {
		if wire == nil {
			continue
		}
		if (wire.FromCircuit == circuit.ID && project.SelectedCircuits[wire.ToCircuit]) || (wire.ToCircuit == circuit.ID && project.SelectedCircuits[wire.FromCircuit]) {
			project.SelectedWires[wire.ID] = true
		}
	}

	project.SelectedCircuits[circuit.ID] = true

	state.SaveProject(ctx, con, project)
	fmt.Printf("\n\nProject: %s\n\n", helpers.SPrettyPrint(project))
	ctx.JSON(http.StatusOK, gin.H{
		"selected-circuits": project.SelectedCircuits,
		"selected-wires":    project.SelectedWires,
	})
}
func SelectCircuitUndo(ctx *gin.Context) {
	var params SelectCircuitParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	var currentScene = project.GetCurrentScene()

	if !project.GetCurrentScene().HasObject(params.CircuitID) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No circuit with provided id",
			"id":    params.CircuitID,
		})
	}

	var circuit = currentScene.GetCircuit(params.CircuitID)

	for _, wire := range circuit.InputWires {
		if wire == nil {
			continue
		}
		delete(project.SelectedWires, wire.ID)
	}
	for _, wire := range circuit.OutputWires {
		if wire == nil {
			continue
		}
		delete(project.SelectedWires, wire.ID)
	}

	delete(project.SelectedCircuits, circuit.ID)

	state.SaveProject(ctx, con, project)
	fmt.Printf("\n\nProject: %s\n\n", helpers.SPrettyPrint(project))
	ctx.JSON(http.StatusOK, gin.H{
		"selected-circuits": project.SelectedCircuits,
		"selected-wires":    project.SelectedWires,
	})

}
