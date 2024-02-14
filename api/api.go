package api

import "github.com/gin-gonic/gin"

var Router *gin.Engine
var Api *gin.RouterGroup

func CreateRouter() {
	Router = gin.Default()
	Api = Router.Group("/api")
}
