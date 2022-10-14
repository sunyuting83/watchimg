package router

import (
	"worldimg/controller"
	utils "worldimg/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY string) *gin.Engine {
	router := gin.Default()
	router.Use(utils.CORSMiddleware())
	api := router.Group("/api")
	api.Use(utils.SetConfigMiddleWare(SECRET_KEY))
	{
		router.GET("/", controller.Index)
		api.POST("/upimg", controller.FileHandler)
		api.GET("/data", utils.VerifyMiddleware(), controller.GetData)
		api.DELETE("/delone", utils.VerifyMiddleware(), controller.DeleteOne)
		api.POST("/login", controller.FileHandler)
		router.Static("/static", "./upimg")
	}

	return router
}
