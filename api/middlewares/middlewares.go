package middlewares

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/abc401/digital-logic-simulator/api/auth"
	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CorsMiddleWare(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", ctx.Request.Header.Get("Origin"))
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, ProjectID, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if ctx.Request.Method == "OPTIONS" {
		fmt.Println("\nOptions request received\n")
		ctx.AbortWithStatus(204)
		return
	}

	ctx.Next()
}

func ProjectIDMiddleWare(ctx *gin.Context) {
	fmt.Println("\nProjectID middleware\n")
	var idStr = ctx.Request.Header.Get("ProjectID")
	var id, err = strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		fmt.Printf("Invalid project id: %d", id)
		ctx.Abort()
		return
	}
	ctx.Set("projectID", id)
	fmt.Printf("Project ID: %d", id)
	ctx.Next()
}

func PrintReqBody(ctx *gin.Context) {
	var body map[string]interface{}

	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil && err != io.EOF {

		fmt.Println("Could not read request body: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not read request body",
		})
		panic("")
	}

	fmt.Println("Request Body: ", helpers.SPrettyPrint(body))

	ctx.Next()
}

func Auth(context *gin.Context) {
	fmt.Println("Auth middleware")
	token := context.Request.Header.Get("Authorization")
	var tokenString string
	if strings.HasPrefix(token, "Bearer ") {
		tokenString = strings.TrimPrefix(token, "Bearer ")
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
		fmt.Print("aborted")
		context.Abort()
		return
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		fmt.Print("aborted")
		context.Abort()
		return
	}
	var userInfo = types.UserInfo{
		UID:   claims.UID,
		Email: claims.Email,
		UName: claims.Username,
	}
	context.Set("userinfo", userInfo)
	context.Next()
}
