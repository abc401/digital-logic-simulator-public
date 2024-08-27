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

func RenameICDo(ctx *gin.Context) {
	type Params struct {
		ID   types.IDType
		From string
		To   string
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
	if _, ok := project.ICs[params.ID]; !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no ic with specified id",
			"id":    params.ID,
		})
		return
	}

	project.ICs[params.ID] = params.To
	project.Scenes[params.ID].Name = params.To

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})

	fmt.Printf("\n\nProject: %s\n\n", helpers.SPrettyPrint(project))
}

func RenameICUndo(ctx *gin.Context) {
	type Params struct {
		ID   types.IDType
		From string
		To   string
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
	if _, ok := project.ICs[params.ID]; !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no ic with specified id",
			"id":    params.ID,
		})
		return
	}

	project.ICs[params.ID] = params.From
	project.Scenes[params.ID].Name = params.From

	state.SaveProject(ctx, con, project)
	ctx.JSON(http.StatusOK, gin.H{})

	fmt.Printf("\n\nProject: %s\n\n", helpers.SPrettyPrint(project))
}
