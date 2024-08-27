package actions

import (
	"log"
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/gin-gonic/gin"
)

func PasteFromClipboardDo(ctx *gin.Context) {
	type Params struct {
		TargetSceneID types.IDType `json:"targetSceneID"`
		Circuits      []Circuit    `binding:"required" json:"circuits"`
		Wires         []Wire       `binding:"required" json:"wires"`
	}
	var params Params

	if !helpers.BindParams(&params, ctx) {
		return
	}
	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	var targetScene, exists = project.Scenes[params.TargetSceneID]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "scene with provided id does not exist",
			"sceneID": params.TargetSceneID,
		})
		return
	}
	for _, circuit := range params.Circuits {
		if targetScene.HasObject(circuit.ID) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":     "provided id for the circuit is already taken",
				"circuitID": circuit.ID,
			})

		}

		for name := range circuit.Props {
			var _, err = state.GetPropSetter(name)
			if err != nil {
				log.Printf("Could not get prop setter with name: %s for %s circuit", name, circuit.CircuitType)
			}

		}
	}
	for _, wire := range params.Wires {
		if targetScene.HasObject(wire.ID) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":  "provided id for the wire is already taken",
				"wireID": wire.ID,
			})
		}
	}

	project.SelectedCircuits = map[types.IDType]bool{}
	project.SelectedWires = map[types.IDType]bool{}

	for _, circuit := range params.Circuits {
		var err = targetScene.AddCircuit(circuit.ID, circuit.ToTypesCircuit())
		if err != nil {
			log.Panicf("circuit with id `%d` shouldn't exist", circuit.ID)
		}
		project.SelectedCircuits[circuit.ID] = true
	}
	for _, wire := range params.Wires {
		var err = targetScene.AddWire(wire.ID, wire.ToTypesWire())
		if err != nil {
			log.Panicf("wire with id `%d` shouldn't exist", wire.ID)
		}
		project.SelectedWires[wire.ID] = true
	}

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})
}

func PasteFromClipboardUndo(ctx *gin.Context) {
	type Params struct {
		SceneID    types.IDType   `json:"targetSceneID"`
		CircuitIDs []types.IDType `binding:"required" json:"circuitIDs"`
		WireIDs    []types.IDType `binding:"required" json:"wireIDs"`
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

	var targetScene, exists = project.Scenes[params.SceneID]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "scene with provided id does not exist",
			"sceneID": params.SceneID,
		})
		return
	}
	for _, id := range params.CircuitIDs {
		if targetScene.GetCircuit(id) == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":     "circuit with provided id does not exist",
				"circuitID": id,
			})
			return
		}

	}
	for _, id := range params.WireIDs {
		if targetScene.GetWire(id) == nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":  "wire with provided id does not exist",
				"wireID": id,
			})
			return
		}

	}

	for _, id := range params.CircuitIDs {
		var err = targetScene.DeleteCircuit(id)
		if err != nil {
			log.Panicf("circuit with id `%d` has to exist", id)
		}
	}
	for _, id := range params.WireIDs {
		var err = targetScene.DeleteWire(id)
		if err != nil {
			log.Panicf("wire with id `%d` has to exist", id)
		}
	}

	project.SelectedCircuits = map[types.IDType]bool{}
	project.SelectedWires = map[types.IDType]bool{}

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})
}
