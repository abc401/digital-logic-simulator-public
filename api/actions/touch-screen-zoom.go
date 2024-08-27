package actions

import (
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/math"
	"github.com/gin-gonic/gin"
)

func TouchScreenZoomDo(ctx *gin.Context) {
	type Params struct {
		EndingView math.ViewManager `binding:"required"`
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

	project.View = params.EndingView

	state.SaveProject(ctx, con, project)

	ctx.JSON(http.StatusOK, gin.H{})
}

func TouchScreenZoomUndo(ctx *gin.Context) {
	type Params struct {
		StartingView math.ViewManager `binding:"required"`
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
	project.View = params.StartingView

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})
}
