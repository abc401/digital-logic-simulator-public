package actions

import (
	"fmt"
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/math"
	"github.com/gin-gonic/gin"
)

type MouseZoomParams struct {
	ZoomOriginScr  math.Vec2 `binding:"required"`
	ZoomLevelDelta float64   `binding:"required"`
}

func MouseZoomDo(ctx *gin.Context) {
	var params MouseZoomParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}

	fmt.Printf("\n\nPrevious View: %s\n", helpers.SPrettyPrint(project.View))
	project.View.MouseZoom(params.ZoomOriginScr, project.View.ZoomLevel+params.ZoomLevelDelta)

	state.SaveProject(ctx, con, project)
	fmt.Printf("Current View: %s\n\n", helpers.SPrettyPrint(project.View))
	ctx.JSON(http.StatusOK, gin.H{
		"current-view": project.View,
	})
}
func MouseZoomUndo(ctx *gin.Context) {
	var params MouseZoomParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}

	fmt.Printf("\n\nPrevious View: %s\n", helpers.SPrettyPrint(project.View))
	project.View.MouseZoom(params.ZoomOriginScr, project.View.ZoomLevel-params.ZoomLevelDelta)

	state.SaveProject(ctx, con, project)
	fmt.Printf("Current View: %s\n\n", helpers.SPrettyPrint(project.View))
	ctx.JSON(http.StatusOK, gin.H{
		"current-view": project.View,
	})
}
