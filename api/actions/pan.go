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

type PanParams struct {
	DeltaScr math.Vec2 `binding:"required"`
}

func PanDo(ctx *gin.Context) {
	var params PanParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	fmt.Printf("\n\nPrevious View: %s\n", helpers.SPrettyPrint(project.View))
	project.View.PanOffset = project.View.PanOffset.Add(params.DeltaScr)

	state.SaveProject(ctx, con, project)
	fmt.Printf("Current View: %s\n\n", helpers.SPrettyPrint(project.View))
	ctx.JSON(http.StatusOK, gin.H{
		"current-view": project.View,
	})
}

func PanUndo(ctx *gin.Context) {
	var params PanParams
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}
	fmt.Printf("Previous View: %s", helpers.SPrettyPrint(project.View))
	project.View.PanOffset = project.View.PanOffset.Sub(params.DeltaScr)

	state.SaveProject(ctx, con, project)
	fmt.Printf("Current View: %s", helpers.SPrettyPrint(project.View))
	ctx.JSON(http.StatusOK, gin.H{
		"current-view": project.View,
	})
}
