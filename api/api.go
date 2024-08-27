package api

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/abc401/digital-logic-simulator/api/actions"
	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/middlewares"
	"github.com/abc401/digital-logic-simulator/api/state"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/abc401/digital-logic-simulator/models"

	"github.com/gin-gonic/gin"
)

var AllowedOrigins = []string{
	"http://localhost5173",
}

func ConfigHandelers(router gin.IRouter) {
	router.Use(middlewares.PrintReqBody)

	// router.GET("/project", GetProjectState)
	// router.GET("/project/reset-state", ResetProjectState)

	router.POST("/signup", SignUp)
	router.POST("/signin", GenToken)
	{
		var create = router.Group("/dbproject-create")
		create.Use(middlewares.Auth)
		create.POST("/create", CreateProject)
	}

	{
		var dbProject = router.Group("/dbproject")
		dbProject.Use(middlewares.Auth)
		dbProject.Use(middlewares.ProjectIDMiddleWare)
		dbProject.GET("/get", GetProject)
		dbProject.GET("/all", ProjectList)
		// dbProject.POST("/select", SelectProject)
	}
	{

		router.GET("/project-list", middlewares.Auth, ProjectList)
	}

	{
		tutorials := router.Group("/tutorials")
		tutorials.GET("/nav", TutorialsNav)
		tutorials.GET("/:link_title", GetTutorial)

		// The `link_title` parameter in this end point probably needs to be renamed to something like `tutorials_id`
		// but I'm getting very bored of this project and its a pain in the ass to do so, and thus, I'm not going to.
		tutorials.GET("/:link_title/quiz", GetQuiz)
	}

	{
		var user = router.Group("/user")
		user.Use(middlewares.Auth)
		user.GET("/info", UserInfo)
		user.GET("/is-signed-in", IsSignedIn)
	}

	{
		action := router.Group("/action")
		actions.ConfigHandlers(action)
	}
}

func IsSignedIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})

}

func UserInfo(ctx *gin.Context) {

	var userinfo, exists = ctx.Get("userinfo")
	if !exists {
		panic("Auth middleware did not provide email")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"userinfo": userinfo,
	})
}

func ResetProjectState(ctx *gin.Context) {
	var con = db.GetGormDBCon()
	*state.GetProject(ctx, con) = state.NewProject()
	ctx.JSON(http.StatusOK, gin.H{})

}

func GetProjectState(ctx *gin.Context) {
	var con = db.GetGormDBCon()
	ctx.JSON(http.StatusOK, state.GetProject(ctx, con))
}

func GetTutorial(ctx *gin.Context) {
	title, found := ctx.Params.Get("link_title")
	if !found {
		panic("Wrong Parameter name")
	}
	fmt.Println("Title: ", title)
	con := db.GetGormDBCon()
	tutorial := models.Article{}

	con.Model(&models.Article{}).Where("link_title = ?", title).Find(&tutorial)
	fmt.Println("Tutorial: ", tutorial.LinkTitle)

	type Res struct {
		ID           uint    `json:"id"`
		LinkTitle    string  `json:"link_title"`
		DisplayTitle string  `json:"display_title"`
		Previous     *string `json:"previous"`
		Next         *string `json:"next"`
		Content      string  `json:"content"`
	}
	res := Res{
		ID:           tutorial.ID,
		LinkTitle:    tutorial.LinkTitle,
		DisplayTitle: tutorial.DisplayTitle,
		Content:      tutorial.Content,
	}
	if tutorial.Previous.Valid {
		res.Previous = &tutorial.Previous.String
	}
	if tutorial.Next.Valid {
		res.Next = &tutorial.Next.String
	}

	fmt.Printf("\n\nTutorial: %s\n\n", helpers.SPrettyPrint(res))

	ctx.JSON(http.StatusOK, res)

}

type Res struct {
	ID           uint    `json:"id"`
	LinkTitle    string  `json:"link_title"`
	DisplayTitle string  `json:"display_title"`
	Previous     *string `json:"previous"`
	Next         *string `json:"next"`
}

func TutorialsNav(ctx *gin.Context) {

	con := db.GetGormDBCon()
	tutorials := []models.Article{}
	con.Select("id", "link_title", "display_title", "previous", "next").Find(&tutorials)
	res := []Res{}
	for i := 0; i < len(tutorials); i++ {
		tutorial := tutorials[i]

		currRes := Res{
			ID:           tutorial.ID,
			LinkTitle:    tutorial.LinkTitle,
			DisplayTitle: tutorial.DisplayTitle,
		}
		if tutorial.Next.Valid {
			currRes.Next = &tutorial.Next.String
		} else {
			currRes.Next = nil
		}
		if tutorial.Previous.Valid {

			currRes.Previous = &tutorial.Previous.String
		} else {
			currRes.Previous = nil
		}

		res = append(res, currRes)
	}
	fmt.Printf("Tutorial Nav: %s", helpers.SPrettyPrint(res))
	ctx.JSON(http.StatusOK, res)
}

func GetQuiz(ctx *gin.Context) {

	// I am getting very bored of this project and its a pain in the ass changing this parameter's name to be something like
	// `tutorials_id` so i'm not going to bother with that for the sake of my sanity
	tutorialID, found := ctx.Params.Get("link_title")

	if !found {
		panic("Wrong Parameter name")
	}
	fmt.Println("Tutorial ID: ", tutorialID)

	con := db.GetGormDBCon()
	var tutorial = Res{}
	mcqs := []*models.MCQ{}

	con.Model(&models.Article{}).Select("id", "link_title", "display_title", "previous", "next").Where("id = ?", tutorialID).Find(&tutorial)
	con.Model(&models.MCQ{}).Where("article_id = ?", tutorialID).Find(&mcqs)
	const n = 5
	for i := 0; uint64(i) < n; i++ {
		var randomIdx = rand.Intn(len(mcqs))
		mcqs[i], mcqs[randomIdx] = mcqs[randomIdx], mcqs[i]
	}

	// fmt.Printf("\n\nTutorial: %s\n\n", helpers.SPrettyPrint(tutorial))
	ctx.JSON(http.StatusOK, gin.H{
		"tutorial": tutorial,
		"mcqs":     mcqs[0:n],
	})

}
