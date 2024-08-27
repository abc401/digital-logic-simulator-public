package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/models"
	"github.com/gin-gonic/gin"
)

func CreateProject(ctx *gin.Context) {
	type Params struct {
		ProjectName string `json:"name" binding:"required"`
	}

	var params Params
	if !helpers.BindParams(&params, ctx) {
		return
	}

	var userInfoAny, exists = ctx.Get("userinfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		return
	}

	var userInfo, ok = userInfoAny.(types.UserInfo)
	if !ok {
		log.Panic("[CreateProject] userinfo was not of correct type")
	}

	var con = db.GetGormDBCon()

	var queryResult = []models.MarshaledProject{}
	con.Where("user_id = ? and name = ?", userInfo.UID, params.ProjectName).First(&queryResult)
	if len(queryResult) != 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":         "project with provided name already exists",
			"provided-name": params.ProjectName,
		})
	}

	var unmarshaledProject = state.NewProject()
	var bytes, err = json.Marshal(unmarshaledProject)
	if err != nil {
		log.Panic("could not marshal the project data into json")
		return
	}

	var marshaledProject = models.MarshaledProject{
		Name:   params.ProjectName,
		Json:   string(bytes),
		UserID: userInfo.UID,
	}

	con.Create(&marshaledProject)
	ctx.JSON(http.StatusOK, gin.H{"id": marshaledProject.ID})
}

func GetProject(ctx *gin.Context) {
	var userInfoAny, exists = ctx.Get("userinfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		return
	}

	var userInfo, ok = userInfoAny.(types.UserInfo)
	if !ok {
		log.Panicf("Userinfo could not be converted to middlewares.UserInfo: %+v", userInfo)
	}

	var projectIDAny any
	projectIDAny, exists = ctx.Get("projectID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "projectID not specified"})
		return
	}
	var projectIDInt int
	projectIDInt, ok = projectIDAny.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id", "projectID": projectIDAny})
		return
	}
	if projectIDInt < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id", "projectID": projectIDAny})
		return
	}
	var projectID = uint(projectIDInt)
	var con = db.GetGormDBCon()
	var project, err = db.GetUnmarshaledProject(con, userInfo.UID, projectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "required project does not exist"})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

func ProjectList(ctx *gin.Context) {
	var userInfoAny, exists = ctx.Get("userinfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		return
	}
	var userInfo, ok = userInfoAny.(types.UserInfo)
	if !ok {
		log.Panicf("Userinfo could not be converted to middlewares.UserInfo: %+v", userInfo)
	}
	var projectMetaData = db.GetAllProjectsMetaData(userInfo.UID)
	ctx.JSON(http.StatusOK, gin.H{"list": projectMetaData})

}
