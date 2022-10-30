package router

import (
	"worldimg/controller"
	utils "worldimg/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string) *gin.Engine {
	router := gin.Default()
	router.Use(utils.CORSMiddleware())
	api := router.Group("/api")
	api.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath))
	{
		router.GET("/", controller.Index)
		api.POST("/upimg", controller.FileHandler)
		api.GET("/data", utils.VerifyMiddleware(), controller.GetData)
		api.GET("/odata", utils.VerifyMiddleware(), controller.GetOldData)
		api.DELETE("/delone", utils.VerifyMiddleware(), controller.DeleteOne)
		api.DELETE("/delist", utils.VerifyMiddleware(), controller.DeleteList)
		api.POST("/login", controller.Sgin)
		router.Static("/static", "./upimg")
	}

	return router
}
