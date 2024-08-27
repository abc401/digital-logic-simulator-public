package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/math"
	"github.com/abc401/digital-logic-simulator/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var project = types.ProjectType{
	Scenes: map[types.IDType]*types.Scene{
		types.DEFAULT_SCENE_ID: {
			ID:        types.DEFAULT_SCENE_ID.ToNullable(),
			Name:      types.DEFAULT_SCENE_NAME,
			ICInputs:  types.NullableID{},
			ICOutputs: types.NullableID{},
			Circuits:  map[types.IDType]*types.Circuit{},
			Wires:     map[types.IDType]*types.Wire{},
		},
	},
	CurrentSceneID:   types.DEFAULT_SCENE_ID,
	SelectedCircuits: map[types.IDType]bool{},
	SelectedWires:    map[types.IDType]bool{},
	View:             math.NewViewManager(),
	ICs:              map[types.IDType]string{},
}

func NewProject() types.ProjectType {
	return types.ProjectType{
		Scenes: map[types.IDType]*types.Scene{
			types.DEFAULT_SCENE_ID: {
				ID:        types.DEFAULT_SCENE_ID.ToNullable(),
				Name:      types.DEFAULT_SCENE_NAME,
				ICInputs:  types.NullableID{},
				ICOutputs: types.NullableID{},
				Circuits:  map[types.IDType]*types.Circuit{},
				Wires:     map[types.IDType]*types.Wire{},
			},
		},
		CurrentSceneID:   types.DEFAULT_SCENE_ID,
		SelectedCircuits: map[types.IDType]bool{},
		SelectedWires:    map[types.IDType]bool{},
		View:             math.NewViewManager(),
		ICs:              map[types.IDType]string{},
	}

}

func GetProject(ctx *gin.Context, con *gorm.DB) *types.ProjectType {
	var userInfoAny, exists = ctx.Get("userinfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		return nil
	}

	var userInfo, ok = userInfoAny.(types.UserInfo)
	if !ok {
		log.Panicf("Userinfo could not be converted to middlewares.UserInfo: %+v", userInfo)
	}

	var projectIDAny any
	projectIDAny, exists = ctx.Get("projectID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "projectID not specified"})
		return nil
	}

	var projectIDInt int
	projectIDInt, ok = projectIDAny.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id", "projectID": projectIDAny})
		return nil
	}
	if projectIDInt < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id", "projectID": projectIDInt})
		return nil
	}
	var projectID = uint(projectIDInt)

	var err error
	var project *types.ProjectType

	project, err = db.GetUnmarshaledProject(con, userInfo.UID, projectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "required project does not exist"})
		return nil
	}

	return project
}

func SaveProject(ctx *gin.Context, con *gorm.DB, project *types.ProjectType) {
	var userInfoAny, exists = ctx.Get("userinfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		log.Panic("impossible!!")
	}

	var userInfo, ok = userInfoAny.(types.UserInfo)
	if !ok {
		log.Panicf("Userinfo could not be converted to middlewares.UserInfo: %+v", userInfo)
	}

	var projectIDAny any
	projectIDAny, exists = ctx.Get("projectID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "projectID not specified"})
		log.Panic("impossible!!")
	}

	var projectIDInt int
	projectIDInt, ok = projectIDAny.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id", "projectID": projectIDAny})
		log.Panic("impossible!!")
	}
	if projectIDInt < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id", "projectID": projectIDInt})
		log.Panic("impossible!!")
	}
	var projectID = uint(projectIDInt)

	var marshaledProject, err = json.Marshal(project)
	if err != nil {
		log.Panic("Could not marshal project")
	}

	con.Model(&models.MarshaledProject{}).Where("user_id = ? and id = ?", userInfo.UID, projectID).Update("json", marshaledProject)
}

var DefaultCircuits = map[string]types.Circuit{
	"Input": {
		CircuitType: "Input",

		NInputPins:  0,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label": "Input",
			"value": false,
		},
	},
	"And": {
		CircuitType: "And",

		NInputPins:  2,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label":  "And",
			"inputs": 2,
		},
	},
	"Or": {
		CircuitType: "Or",

		NInputPins:  2,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label":  "Or",
			"inputs": 2,
		},
	},
	"Not": {
		CircuitType: "Not",

		NInputPins:  1,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label": "Not",
		},
	},
	"Nand": {
		CircuitType: "Nand",

		NInputPins:  2,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label":  "Nand",
			"inputs": 2,
		},
	},
	"Nor": {
		CircuitType: "Nor",

		NInputPins:  2,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label":  "Nor",
			"inputs": 2,
		},
	},
	"Xor": {
		CircuitType: "Xor",

		NInputPins:  2,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label":  "Xor",
			"inputs": 2,
		},
	},
	"Xnor": {
		CircuitType: "Xnor",

		NInputPins:  2,
		NOutputPins: 1,
		Props: types.CircuitProps{
			"label":  "Xnor",
			"inputs": 2,
		},
	},
}

func NewICInputs() *types.Circuit {
	return &types.Circuit{
		ID:          CUSTOM_CIRCUIT_INPUTS_ID,
		CircuitType: "customcircuitinputs",
		PosWrl:      math.NewVec2(90, 220),
		NInputPins:  0,
		NOutputPins: 1,
		InputWires:  []*types.Wire{},
		OutputWires: make([]*types.Wire, 1),
		Props: types.CircuitProps{
			"label": "CustomCircuitInputs",
		},
	}
}

func NewICOutputs() *types.Circuit {
	return &types.Circuit{
		ID:          CUSTOM_CIRCUIT_OUTPUTS_ID,
		CircuitType: "customcircuitoutputs",
		PosWrl:      math.NewVec2(90, 220),
		NInputPins:  1,
		NOutputPins: 0,
		InputWires:  make([]*types.Wire, 1),
		OutputWires: []*types.Wire{},
		Props: types.CircuitProps{
			"label": "CustomCircuitOutputs",
		},
	}
}

func SetCircuitLabel(circuit *types.Circuit, value string) bool {
	if value == "" {
		return false
	}
	var isWhiteSpace = true
	for _, char := range value {
		if !unicode.IsSpace(char) {
			isWhiteSpace = false
		}
	}
	if isWhiteSpace {
		return false
	}

	circuit.Props["label"] = value
	return true
}

func SetInputValue(circuit *types.Circuit, value string) bool {

	var valueBool, err = strconv.ParseBool(strings.Trim(value, " \t"))
	if err != nil {
		return false
	}

	if circuit.CircuitType != "input" {

		str, err := json.MarshalIndent(circuit, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("\n\n[Error] Tried to set 'value' prop of: %s\n\n", str)

	}
	circuit.Props["value"] = valueBool
	return true
}

func SetCircuitInputs(circuit *types.Circuit, value string) bool {
	if _, ok := circuit.Props["inputs"]; !ok {
		str, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("\n\n[Error] Tried to set 'inputs' prop of circuit to: %s\n\n", str)

		str, err = json.MarshalIndent(circuit, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Circuit: %s", str)
	}

	var nConsumerPins, err = strconv.ParseInt(value, 10, 64)
	if err != nil || nConsumerPins < 0 {
		str, err := json.MarshalIndent(value, "", "	")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("[Info] Error parsing int from string: %s", str)
		return false
	}

	if len(circuit.InputWires) == int(nConsumerPins) {
		fmt.Println("[Info] len(circuit.InputWires) == int(nConsumerPins)")
		return true
	}

	var nConnectedPins = uint64(0)

	for _, wire := range circuit.InputWires {
		if wire != nil {
			nConnectedPins++
		}
	}
	if nConnectedPins > uint64(nConsumerPins) {
		fmt.Println("[Info] nConnectedPins > uint64(nConsumerPins)")
		return false
	}

	var newPins = make([]*types.Wire, nConsumerPins)

	var pinIndex = uint64(0)
	for _, pin := range circuit.InputWires {
		if pin != nil {
			newPins[pinIndex] = pin
			pin.ToPin = pinIndex
			pinIndex++
		}

	}
	circuit.InputWires = newPins
	circuit.Props["inputs"] = nConsumerPins
	circuit.NInputPins = uint64(nConsumerPins)

	fmt.Printf("[Info] Successfully set value to %v", nConsumerPins)
	return true
}

type PropSetter func(*types.Circuit, string) bool

var commonPropSetters = map[string]PropSetter{
	"inputs": SetCircuitInputs,
}

var CircuitPropSetters = map[string]map[string]PropSetter{
	"*": {
		"label": SetCircuitLabel,
	},
	"input": {
		"value": SetInputValue,
	},
	"and":  commonPropSetters,
	"or":   commonPropSetters,
	"not":  commonPropSetters,
	"nand": commonPropSetters,
	"nor":  commonPropSetters,
}

func GetPropSetter(name string) (*PropSetter, error) {
	for _, props := range CircuitPropSetters {
		var propSetter, ok = props[name]
		if ok {
			return &propSetter, nil
		}
	}
	return nil, errors.New("propSetter doesn't exist")
}

const CUSTOM_CIRCUIT_INPUTS_ID = types.IDType(0)
const CUSTOM_CIRCUIT_OUTPUTS_ID = types.IDType(1)

func NewSceneWithIO(id types.IDType, name string) *types.Scene {
	return &types.Scene{
		ID:        id.ToNullable(),
		Name:      name,
		ICInputs:  CUSTOM_CIRCUIT_INPUTS_ID.ToNullable(),
		ICOutputs: CUSTOM_CIRCUIT_OUTPUTS_ID.ToNullable(),
		Circuits: map[types.IDType]*types.Circuit{
			CUSTOM_CIRCUIT_INPUTS_ID:  NewICInputs(),
			CUSTOM_CIRCUIT_OUTPUTS_ID: NewICOutputs(),
		},
		Wires: map[types.IDType]*types.Wire{},
	}
}
