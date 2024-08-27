package actions

import (
	"github.com/abc401/digital-logic-simulator/api/middlewares"
	"github.com/gin-gonic/gin"
)

func ConfigHandlers(router gin.IRouter) {
	router.Use(middlewares.Auth)
	router.Use(middlewares.ProjectIDMiddleWare)

	router.POST("/pan/do", PanDo)
	router.POST("/pan/undo", PanUndo)

	router.POST("/mouse-zoom/do", MouseZoomDo)
	router.POST("/mouse-zoom/undo", MouseZoomUndo)

	router.POST("/touch-screen-zoom/do", TouchScreenZoomDo)
	router.POST("/touch-screen-zoom/undo", TouchScreenZoomUndo)

	router.POST("/select-circuit/do", SelectCircuitDo)
	router.POST("/select-circuit/undo", SelectCircuitUndo)

	router.POST("/deselect-circuit/do", SelectCircuitUndo)
	router.POST("/deselect-circuit/undo", SelectCircuitDo)

	router.POST("/deselect-all/do", DeselectAllDo)
	router.POST("/deselect-all/undo", DeselectAllUndo)

	router.POST("/drag-selection/do", DragSelectionDo)
	router.POST("/drag-selection/undo", DragSelectionUndo)

	router.POST("/create-circuit/do", CreateCircuitDo)
	router.POST("/create-circuit/undo", CreateCircuitUndo)

	router.POST("/create-wire/do", CreateWireDo)
	router.POST("/create-wire/undo", CreateWireUndo)

	router.POST("/delete-wire/do", CreateWireUndo)
	router.POST("/delete-wire/undo", CreateWireDo)

	router.POST("/create-ic/do", CreateICDo)
	router.POST("/create-ic/undo", CreateICUndo)

	router.POST("switch-scene/do", SwitchSceneDo)
	router.POST("switch-scene/undo", SwitchSceneUndo)

	router.POST("rename-ic/do", RenameICDo)
	router.POST("rename-ic/undo", RenameICUndo)

	router.POST("set-prop/do", SetCircuitPropDo)
	router.POST("set-prop/undo", SetCircuitPropUndo)

	router.POST("delete-selected/do", DeleteSelectedDo)
	router.POST("delete-selected/undo", DeleteSelectedUndo)

	router.POST("paste-from-clipboard/do", PasteFromClipboardDo)
	router.POST("paste-from-clipboard/undo", PasteFromClipboardUndo)

	router.POST("/noop", Noop)
}
