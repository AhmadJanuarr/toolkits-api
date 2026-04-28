package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"toolkits/internal/config"
	"toolkits/internal/services"
	"toolkits/internal/utils"

	"github.com/gin-gonic/gin"
)

type DownloaderHandler struct {
	config    *config.Config
	semaphore chan struct{}
}

func NewDownloaderHandler(cfg *config.Config, sem chan struct{}) *DownloaderHandler {
	return &DownloaderHandler{
		config:    cfg,
		semaphore: sem,
	}
}

func (h *DownloaderHandler) Downloader(c *gin.Context) {
	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		inputURL := c.PostForm("inputURL")
		formatID := c.PostForm("format_id")
		if inputURL == "" || formatID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "inputURL dan FormatID tidak boleh kosong",
			})
			return
		}
		cacheKey := inputURL + "_" + formatID
		if cachedPath, exists := utils.FileCache.Get(cacheKey); exists {
			fileName := filepath.Base(cachedPath)
			c.FileAttachment(cachedPath, fileName)
			return
		}
		resultPath, err := services.ProsessDownload(inputURL, formatID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Gagal mengunduh",
				"error":   err.Error(),
			})
			return
		}
		utils.FileCache.Set(cacheKey, resultPath)
		fileName := filepath.Base(resultPath)
		c.FileAttachment(resultPath, fileName)

	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}
}

func (h *DownloaderHandler) DownloaderGetInfo(c *gin.Context) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()
		//
		inputURL := c.PostForm("inputURL")
		if strings.Contains(inputURL, "/photo/") {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Maaf, untuk saat ini format tiktok berupa slide foto belum mendukung",
			})
			return
		}

		if inputURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "inputURL tidak boleh kosong atau tidak valid",
			})
			return
		}

		info, err := services.ProcessGetInfo(inputURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Gagal mengambil informasi video",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Informasi video berhasil diambil",
			"data": gin.H{
				"info": info,
			},
		})
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}

}
