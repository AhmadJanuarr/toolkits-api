package main

import (
	"log"
	"os"
	"toolkits/internal/config"
	"toolkits/internal/jobs"
	"toolkits/internal/routes"

	"github.com/joho/godotenv"
)

func main() {

	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cookieData := os.Getenv("YOUTUBE_COOKIE_DATA")
	if cookieData != "" {
		os.WriteFile("youtube_cookies.txt", []byte(cookieData), 0644)
	}

	cookieData = os.Getenv("INSTAGRAM_COOKIE_DATA")
	if cookieData != "" {
		os.WriteFile("ig-cookies.txt", []byte(cookieData), 0644)
	}

	cfg := config.LoadConfig()

	cleaner := jobs.StorageCleanup(
		[]string{
			cfg.Storage.ProcessedDir,
			cfg.Storage.CompressedDir,
			cfg.Storage.ResizedDir,
			cfg.Storage.DownloadsDir,
		},
		cfg.Worker.MaxFileAge,
	)
	cleaner.Start()

	router := routes.Route(cfg)

	router.MaxMultipartMemory = cfg.Server.MaxMultipartMemory

	log.Println("Server running on port:", cfg.Server.Port)
	router.Run(":" + cfg.Server.Port)
}
