package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mangalores/case-studies-voiceline/internal/app"
	"github.com/mangalores/case-studies-voiceline/internal/http/handlers"
	"github.com/mangalores/case-studies-voiceline/internal/http/middleware"
)

func NewRouter(uploader *app.UploadService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())

	recordingHandler := handlers.NewRecordingHandler(uploader)
	router.POST("/recordings", recordingHandler.Create)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}
