package actions

import (
	"log"
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/math"
	"github.com/gin-gonic/gin"
)

func DeleteSelectedDo(ctx *gin.Context) {
	type Params struct {
		SceneID    types.IDType   `json:"sceneID"`
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

type Wire struct {
	ID types.IDType `json:"id"`

	FromCircuit types.IDType `json:"fromCircuit"`
	FromPin     uint64       `json:"fromPin"`

	ToCircuit types.IDType `json:"toCircuit"`
	ToPin     uint64       `json:"toPin"`
}

func (wire *Wire) ToTypesWire() *types.Wire {
	return &types.Wire{
		ID:          wire.ID,
		FromCircuit: wire.FromCircuit,
		FromPin:     wire.FromPin,
		ToCircuit:   wire.ToCircuit,
		ToPin:       wire.ToPin,
	}

}

type Circuit struct {
	ID          types.IDType       `json:"id"`
	CircuitType string             `binding:"required" json:"circuitType"`
	PosWrl      math.Vec2          `binding:"required" json:"posWrl"`
	NInputPins  uint64             `json:"nInputPins"`
	NOutputPins uint64             `json:"nOutputPins"`
	InputWires  []*Wire            `json:"inputWires"`
	OutputWires []*Wire            `json:"outputWires"`
	Props       types.CircuitProps `json:"props"`
}

func (circuit *Circuit) ToTypesCircuit() *types.Circuit {
	var inputWires []*types.Wire = []*types.Wire{}
	for _, wire := range circuit.InputWires {
		inputWires = append(inputWires, wire.ToTypesWire())
	}
	var outputWires []*types.Wire = []*types.Wire{}
	for _, wire := range circuit.OutputWires {
		outputWires = append(outputWires, wire.ToTypesWire())
	}

	return &types.Circuit{
		ID:          circuit.ID,
		CircuitType: circuit.CircuitType,
		PosWrl:      circuit.PosWrl,
		NInputPins:  circuit.NInputPins,
		NOutputPins: circuit.NOutputPins,
		InputWires:  inputWires,
		OutputWires: outputWires,
		Props:       circuit.Props,
	}
}

func DeleteSelectedUndo(ctx *gin.Context) {

	type Params struct {
		SceneID         types.IDType `json:"sceneID"`
		Circuits        []Circuit    `binding:"required" json:"circuits"`
		SelectedWires   []Wire       `binding:"required" json:"selectedWires"`
		UnSelectedWires []Wire       `binding:"required" json:"unSelectedWires"`
	}
	var params Params

	if !helpers.BindParams(&params, ctx) {
		return
	}
	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	var targetScene, exists = project.Scenes[params.SceneID]
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "scene with provided id does not exist",
			"sceneID": params.SceneID,
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
	for _, wire := range params.SelectedWires {
		if targetScene.HasObject(wire.ID) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":  "provided id for the wire is already taken",
				"wireID": wire.ID,
			})
		}
	}
	for _, wire := range params.UnSelectedWires {
		if targetScene.HasObject(wire.ID) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":  "provided id for the wire is already taken",
				"wireID": wire.ID,
			})
		}
	}

	for _, circuit := range params.Circuits {
		var err = targetScene.AddCircuit(circuit.ID, circuit.ToTypesCircuit())
		if err != nil {
			log.Panicf("circuit with id `%d` shouldn't exist", circuit.ID)
		}
		project.SelectedCircuits[circuit.ID] = true
	}
	for _, wire := range params.SelectedWires {
		var err = targetScene.AddWire(wire.ID, wire.ToTypesWire())
		if err != nil {
			log.Panicf("wire with id `%d` shouldn't exist", wire.ID)
		}
		project.SelectedWires[wire.ID] = true
	}
	for _, wire := range params.UnSelectedWires {
		var err = targetScene.AddWire(wire.ID, wire.ToTypesWire())
		if err != nil {
			log.Panicf("wire with id `%d` shouldn't exist", wire.ID)
		}
	}

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})
}
