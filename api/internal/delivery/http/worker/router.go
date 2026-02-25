package worker

import (
	"github.com/gin-gonic/gin"
	"github.com/adityawiryaa/api/internal/middleware"
)

func SetupRouter(handler *Handler, apiKey string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogging())

	r.POST("/config", handler.ReceiveConfig)
	r.GET("/config", handler.GetCurrentConfig)
	r.GET("/hit", handler.ExecuteHit)
	r.GET("/hit/:taskId", handler.GetHitResult)

	return r
}
