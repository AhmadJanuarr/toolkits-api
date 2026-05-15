package routes

import (
	"net/http"
	"toolkits/internal/config"
	"toolkits/internal/handlers"

	"toolkits/internal/middlewares"
)

func Route(cfg *config.Config) http.Handler {

	mux := http.NewServeMux()

	globalLimit := make(chan struct{}, cfg.Server.MaxGlobalConcurrent)
	imgHandler := handlers.NewImageHandler(cfg, globalLimit)
	downloaderHandler := handlers.NewDownloaderHandler(cfg, globalLimit)

	mux.HandleFunc("POST /api/v1/image/convert", imgHandler.ConvertImage)
	mux.HandleFunc("POST /api/v1/image/compress-image", imgHandler.CompressionImage)
	mux.HandleFunc("POST /api/v1/image/resize-image", imgHandler.ResizeImage)

	proxyHandler := handlers.NewProxyHandler()

	mux.HandleFunc("POST /api/v1/downloader/info", downloaderHandler.DownloaderGetInfo)
	mux.HandleFunc("GET /api/v1/downloader/download", downloaderHandler.Downloader)
	mux.HandleFunc("GET /api/v1/proxy/image", proxyHandler.ProxyImage)

	var handler http.Handler = mux
	handler = middlewares.CORSMiddleware(cfg)(handler)
	handler = middlewares.RateLimitMiddleware(cfg.Server.RateLimitRPS, cfg.Server.RateLimitBurst)(handler)
	handler = middlewares.TimeoutMiddleware(cfg.Server.ReadTimeout)(handler)
	handler = middlewares.LoggerMiddleware()(handler)
	return handler
}
