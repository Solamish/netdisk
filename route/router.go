package route

import (
	"github.com/gin-gonic/gin"
	"netdisk/handler"
)

func LOAD(router *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {



	//用户相关api
	userGroup := router.Group("/user")
	{
		userGroup.POST("/signup",handler.DoSignUpHandler)
		userGroup.POST("/signin",handler.DoSignInHandler)
	}


	router.Use(handler.HTTPInterceptor())
	fileGroup := router.Group("/file")
	{

		fileGroup.POST("/upload",handler.DoUploadHandler)
		fileGroup.GET("/download",handler.DownloadHandler)
		fileGroup.POST("/fast",handler.TryFastUploadHandler)
		fileGroup.GET("/query",handler.FileQueryHandler)
	}

	mpGroup := router.Group("/multiPart")
	{
		mpGroup.POST("/init",handler.InitialMultipartHandler)
		mpGroup.POST("/upload",handler.UploadPartHandler)
		mpGroup.POST("/complete",handler.CompleteUploadHandler)
	}


	return router
}