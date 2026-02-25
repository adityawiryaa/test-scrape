package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/adityawiryaa/api/internal/middleware"
)

func SetupRouter(handler *Handler, apiKey string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogging())

	protected := r.Group("")
	protected.Use(middleware.APIKeyAuth(apiKey))
	{
		protected.POST("/register", handler.RegisterAgent)
		protected.POST("/config", handler.UpdateConfig)
		protected.GET("/config", handler.GetConfig)
		protected.GET("/config/:version", handler.GetConfigByVersion)
	}

	return r
}
