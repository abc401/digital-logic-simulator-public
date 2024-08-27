package main

import (
	"fmt"
	"os"

	"github.com/abc401/digital-logic-simulator/api"
	"github.com/abc401/digital-logic-simulator/api/middlewares"
	"github.com/abc401/digital-logic-simulator/db"
	"github.com/gin-gonic/gin"
)

func RunServer() {

	router := gin.Default()
	router.Use(middlewares.CorsMiddleWare)

	var apiRouter = router.Group("/api")
	api.ConfigHandelers(apiRouter)
	router.Run()
}

func main() {
	if len(os.Args) == 1 {
		RunServer()
	} else if len(os.Args) == 2 && os.Args[1] == "--configdb" {
		db.AutoMigrate()
		db.AddTutorials(Tutorials)
	} else {
		fmt.Println("Incorrect usage!!")
		fmt.Println("Correct Usage: \n\t go run . [--config]")
		os.Exit(1)
	}
}
