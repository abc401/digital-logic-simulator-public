package actions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/gin-gonic/gin"
)

type CreateICParams struct {
	SceneID types.IDType
	ICName  string `binding:"required"`
}

func CreateICDo(ctx *gin.Context) {
	var params CreateICParams

	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()

	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}

	if params.SceneID == types.DEFAULT_SCENE_ID {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "cannot reference scene 0 (Main scene) for custom circuit",
		})
	}

	if _, ok := project.ICs[params.SceneID]; ok {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "id is already taken",
			"id":    params.SceneID,
		})
	}

	for _, name := range project.ICs {
		if strings.EqualFold(name, params.ICName) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "proposed name is already taken",
				"name":  params.ICName,
			})
		}
	}

	fmt.Printf("id: %+v, idNullable: %+v", params.SceneID, params.SceneID.ToNullable())
	var newScene = state.NewSceneWithIO(params.SceneID, params.ICName)
	project.Scenes[params.SceneID] = newScene
	project.ICs[params.SceneID] = params.ICName

	state.SaveProject(ctx, con, project)
	fmt.Printf("\n\nProject: %s\n\n", helpers.SPrettyPrint(project))
	ctx.JSON(http.StatusOK, gin.H{})
}

func CreateICUndo(ctx *gin.Context) {
	var params CreateICParams

	if !helpers.BindParams(&params, ctx) {
		return
	}

	var con = db.GetGormDBCon()
	var project = state.GetProject(ctx, con)
	if project == nil {
		return
	}

	if params.SceneID == types.DEFAULT_SCENE_ID {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "scene 0 (Main Scene) is not a custom ic",
		})
	}

	if _, ok := project.ICs[params.SceneID]; !ok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no custom ic with given id",
			"id":    params.SceneID,
		})
	}

	delete(project.Scenes, params.SceneID)
	delete(project.ICs, params.SceneID)

	state.SaveProject(ctx, con, project)
	fmt.Printf("\n\nProject: %s\n\n", helpers.SPrettyPrint(project))

	ctx.JSON(http.StatusOK, gin.H{})
}
