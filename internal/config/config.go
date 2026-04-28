package config

import (
	"strings"
	"time"
	"toolkits/internal/utils"
)

type Config struct {
	Server  ServerConfig
	Image   ImageConfig
	Storage StorageConfig
	Worker  WorkerConfig
}
type ServerConfig struct {
	Port                string
	MaxMultipartMemory  int64
	MaxGlobalConcurrent int
	AllowedOrigins      string
	AllowedMethods      []string
	AllowedHeaders      []string
	ExposeHeaders       []string
	RateLimitRPS        float64
	RateLimitBurst      int
}
type ImageConfig struct {
	MaxFileSize  int64
	MaxDimension int

	AllowedFormats []string
	DefaultQuality int
	MinQuality     int
	MaxQuality     int
}
type StorageConfig struct {
	TempDir       string
	UploadDir     string
	ProcessedDir  string
	CompressedDir string
	ResizedDir    string
	DownloadsDir  string
}
type WorkerConfig struct {
	CleanupInterval time.Duration
	MaxFileAge      time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:                utils.GetEnv("PORT", "8080"),
			MaxMultipartMemory:  utils.GetInt64("MAX_MULTIPART_MEMORY", 20<<20),
			AllowedOrigins:      utils.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
			AllowedMethods:      strings.Split(utils.GetEnv("ALLOWED_METHODS", "GET, POST, PUT, DELETE, OPTIONS"), ","),
			AllowedHeaders:      strings.Split(utils.GetEnv("ALLOWED_HEADERS", "Origin, Content-Type, Authorization"), ","),
			ExposeHeaders:       strings.Split(utils.GetEnv("EXPOSE_HEADERS", "Content-Length, Content-Type, Content-Disposition"), ","),
			MaxGlobalConcurrent: utils.GetInt("MAX_GLOBAL_CONCURRENT", 4),
			RateLimitRPS:        2.0,
			RateLimitBurst:      5,
		},
		Image: ImageConfig{
			MaxFileSize:    utils.GetInt64("MAX_FILE_SIZE", 5<<20),
			MaxDimension:   utils.GetInt("MAX_DIMENSION", 4096),
			AllowedFormats: []string{"jpg", "jpeg", "png", "webp"},
			DefaultQuality: utils.GetInt("DEFAULT_QUALITY", 80),
			MinQuality:     utils.GetInt("MIN_QUALITY", 1),
			MaxQuality:     utils.GetInt("MAX_QUALITY", 100),
		},
		Storage: StorageConfig{
			TempDir:       utils.GetEnv("TEMP_DIR", "temp"),
			UploadDir:     utils.GetEnv("UPLOAD_DIR", "temp/uploads"),
			ProcessedDir:  utils.GetEnv("PROCESSED_DIR", "temp/processed"),
			CompressedDir: utils.GetEnv("COMPRESSED_DIR", "temp/compressed"),
			ResizedDir:    utils.GetEnv("RESIZED_DIR", "temp/resized"),
			DownloadsDir:  utils.GetEnv("DOWNLOADS_DIR", "temp/downloads"),
		},
		Worker: WorkerConfig{
			CleanupInterval: utils.GetDuration("CLEANUP_INTERVAL", 5*time.Second),
			MaxFileAge:      utils.GetDuration("MAX_FILE_AGE", 24*time.Hour),
		},
	}
}
