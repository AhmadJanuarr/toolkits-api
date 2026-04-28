package routes

import (
	"toolkits/internal/config"
	"toolkits/internal/handlers"

	"toolkits/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Route(cfg *config.Config) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS CONFIGURATION

	origins := cfg.Server.AllowedOrigins

	config := cors.Config{
		AllowOrigins:     []string{origins},
		AllowMethods:     cfg.Server.AllowedMethods,
		AllowHeaders:     cfg.Server.AllowedHeaders,
		ExposeHeaders:    cfg.Server.ExposeHeaders,
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	router.Use(middlewares.RateLimitMiddleware(cfg.Server.RateLimitRPS, cfg.Server.RateLimitBurst))
	globalLimit := make(chan struct{}, cfg.Server.MaxGlobalConcurrent)

	imgHandler := handlers.NewImageHandler(cfg, globalLimit)
	downloaderHandler := handlers.NewDownloaderHandler(cfg, globalLimit)
	// ROUTES

	v1 := router.Group("/api/v1")

	// Image endpoints
	v1.POST("/image/convert", imgHandler.ConvertImage)
	v1.POST("/image/compress-image", imgHandler.CompressionImage)
	v1.POST("/image/resize-image", imgHandler.ResizeImage)

	// Youtube
	v1.POST("/downloader/info", downloaderHandler.DownloaderGetInfo)
	v1.POST("/downloader/download", downloaderHandler.Downloader)

	return router
}
