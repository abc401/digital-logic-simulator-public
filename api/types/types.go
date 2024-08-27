package types

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/abc401/digital-logic-simulator/math"
)

const DEFAULT_SCENE_NAME = "Main"
const DEFAULT_SCENE_ID = IDType(0)

type UserInfo struct {
	UID   uint
	Email string
	UName string
}

type IDType uint64

func (id IDType) ToNullable() NullableID {
	return NullableID{
		ID:    id,
		Valid: true,
	}
}

type NullableID struct {
	ID    IDType
	Valid bool
}

func (id *NullableID) Some(validID IDType) *NullableID {
	id.ID = validID
	id.Valid = true
	return id
}

func (id *NullableID) None() *NullableID {
	id.Valid = false
	return id
}

func (id *NullableID) Unwrap() (IDType, bool) {
	return id.ID, id.Valid
}

type Scene struct {
	ID        NullableID          `json:"ID"`
	Name      string              `json:"Name"`
	ICInputs  NullableID          `json:"ICInputs"`
	ICOutputs NullableID          `json:"ICOutputs"`
	Circuits  map[IDType]*Circuit `json:"Circuits"`
	Wires     map[IDType]*Wire    `json:"Wires"`
}

func (scene *Scene) HasCircuit(id IDType) bool {
	_, circuitExists := scene.Circuits[id]
	return circuitExists
}
func (scene *Scene) HasObject(id IDType) bool {
	return scene.HasCircuit(id) || scene.HasWire(id)
}

func (scene *Scene) HasWire(id IDType) bool {
	_, wireExists := scene.Wires[id]
	return wireExists
}

func (scene *Scene) GetCircuit(id IDType) *Circuit {
	circuit, ok := scene.Circuits[id]
	if ok {
		return circuit
	}
	return nil
}
func (scene *Scene) GetWire(id IDType) *Wire {
	wire, ok := scene.Wires[id]
	if ok {
		return wire
	}
	return nil
}

func (scene *Scene) AddCircuit(id IDType, circuit *Circuit) error {
	if scene.HasObject(id) {
		return errors.New("id is already taken")
	}

	scene.Circuits[id] = circuit

	return nil
}

func (scene *Scene) AddWire(id IDType, wire *Wire) error {
	if scene.HasObject(id) {
		return errors.New("id is already taken")
	}

	scene.Wires[id] = wire

	return nil
}

// Returns error if there is no circuit registered
// with the id `id`. The wires connected with the
// circuit being deleted have to be manually deleted
// using the `DeleteWire` function.
func (scene *Scene) DeleteCircuit(id IDType) error {
	if !scene.HasCircuit(id) {
		return errors.New("circuit does not exist")
	}

	delete(scene.Circuits, id)

	return nil
}

// Returns error if there is no wire registered
// with the id `id`. The caller has to ensure that the wire
// being deleted is not connected to any circuit.
func (scene *Scene) DeleteWire(id IDType) error {
	if !scene.HasWire(id) {
		return errors.New("wire does not exist")
	}

	delete(scene.Wires, id)

	return nil

}

func (project *ProjectType) DeselectAll() {
	project.SelectedCircuits = map[IDType]bool{}
	project.SelectedWires = map[IDType]bool{}
}

type ProjectType struct {
	Scenes           map[IDType]*Scene `json:"Scenes"`
	View             math.ViewManager  `json:"View"`
	CurrentSceneID   IDType            `json:"CurrentSceneID"`
	SelectedCircuits map[IDType]bool   `json:"SelectedCircuits"`
	SelectedWires    map[IDType]bool   `json:"SelectedWires"`
	ICs              map[IDType]string `json:"ICs"`
}

func (project *ProjectType) ReEvaluateICs(targetSceneID IDType, icKindID IDType) {

	if icKindID == DEFAULT_SCENE_ID {
		return
	}
	var targetScene, ok = project.Scenes[targetSceneID]
	if !ok {
		log.Panicf("[ReEvaluateICs] attempted to reevaluate scene that does not exist. scene id: %v", targetSceneID)
	}

	var icKind string
	icKind, ok = project.ICs[icKindID]
	if !ok {
		log.Panicf("[ReEvaluateICs] attempted to reevaluate ic type that doesn't exist. ickind id: %v", icKindID)
	}
	var icScene *Scene
	icScene, ok = project.Scenes[icKindID]
	if !ok {
		log.Panicf("[ReEvaluateICs] entry for ic is present in project.ics but not in project.scenes ickind id: %v", icKindID)
	}

	if !icScene.ICInputs.Valid || !icScene.ICOutputs.Valid {
		log.Panicf("[ReEvaluateICs] ic scene doesn't have icinputs or icoutputs. ickind id: %v", icKindID)
	}

	var icInputs = icScene.GetCircuit(icScene.ICInputs.ID)
	if icInputs == nil {
		log.Panicf("[ReEvaluateICs] id is valid but icInputs not present. ickind id: %v", icKindID)
	}

	var icOutputs = icScene.GetCircuit(icScene.ICOutputs.ID)
	if icOutputs == nil {
		log.Panicf("[ReEvaluateICs] id is valid but icOutputs not present. ickind id: %v", icKindID)
	}

	fmt.Printf("Reevaluating scene: %v for ictypeid: %v", targetScene.Name, icKind)

	for _, circuit := range targetScene.Circuits {
		if strings.EqualFold(circuit.CircuitType, icKind) {
			circuit.NInputPins = icInputs.NOutputPins - 1
			circuit.NOutputPins = icOutputs.NInputPins - 1
			for int(circuit.NInputPins) > len(circuit.InputWires) {
				circuit.InputWires = append(circuit.InputWires, nil)
			}
			for int(circuit.NOutputPins) > len(circuit.OutputWires) {
				circuit.OutputWires = append(circuit.OutputWires, nil)
			}
		}

	}

}

func (project *ProjectType) GetCurrentScene() *Scene {
	return project.Scenes[project.CurrentSceneID]
}

type CircuitProps map[string]interface{}

type Circuit struct {
	ID          IDType
	CircuitType string
	PosWrl      math.Vec2
	NInputPins  uint64
	NOutputPins uint64
	InputWires  []*Wire
	OutputWires []*Wire
	Props       CircuitProps
}

type Wire struct {
	ID IDType

	FromCircuit IDType
	FromPin     uint64

	ToCircuit IDType
	ToPin     uint64
}
