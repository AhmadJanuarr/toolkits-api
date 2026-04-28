package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"toolkits/internal/config"
	"toolkits/internal/services"
	"toolkits/internal/utils"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	config    *config.Config
	semaphore chan struct{}
}

func NewImageHandler(cfg *config.Config, sem chan struct{}) *ImageHandler {
	return &ImageHandler{
		config:    cfg,
		semaphore: sem,
	}
}

func (h *ImageHandler) ConvertImage(c *gin.Context) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		targetFormat := strings.ToLower(c.PostForm("format"))
		if targetFormat == "" || !utils.Contains(h.config.Image.AllowedFormats, targetFormat) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Format target tidak valid atau tidak didukung (gunakan: jpg, jpeg, png, webp)",
			})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "File tidak ditemukan atau tidak valid",
				"error":   err.Error(),
			})
			return
		}

		if err := utils.ValidateImageFile(file); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		if file.Size > h.config.Image.MaxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Ukuran file tidak boleh lebih dari %d MB", h.config.Image.MaxFileSize>>20),
			})
			return
		}

		safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
		srcPath := filepath.Join(h.config.Storage.TempDir, safeFilename)

		if err := c.SaveUploadedFile(file, srcPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Gagal menyimpan file sementara",
				"error":   err.Error(),
			})
			return
		}
		defer os.Remove(srcPath)

		// process image conversion
		resultPath, err := services.ProcessImageConversion(srcPath, targetFormat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Gagal mengonversi gambar",
				"error":   err.Error(),
			})
			return
		}

		c.File(resultPath)

	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}
}

func (h *ImageHandler) CompressionImage(c *gin.Context) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "File tidak ditemukan atau tidak valid",
				"error":   err.Error(),
			})
			return
		}

		qualityStr := c.PostForm("quality")
		if qualityStr == "" {
			qualityStr = "80"
		}

		quality, err := strconv.Atoi(qualityStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Quality harus berupa angka",
				"error":   err.Error(),
			})
			return
		}
		if quality < h.config.Image.MinQuality || quality > h.config.Image.MaxQuality {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Quality harus antara %d dan %d", h.config.Image.MinQuality, h.config.Image.MaxQuality),
			})
			return
		}

		safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
		srcPath := filepath.Join(h.config.Storage.UploadDir, safeFilename)

		if err := c.SaveUploadedFile(file, srcPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Gagal menyimpan file sementara",
				"error":   err.Error(),
			})
			return
		}
		defer os.Remove(srcPath)

		resultPath, err := services.ProcessImageCompression(srcPath, quality)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Gagal mengompres gambar",
				"error":   err.Error(),
			})
			return
		}

		c.File(resultPath)

	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}
}

func (h *ImageHandler) ResizeImage(c *gin.Context) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "File wajib diupload", "error": err.Error()})
			return
		}

		widthStr := c.PostForm("width")
		heightStr := c.PostForm("height")

		width, errW := strconv.Atoi(widthStr)
		height, errH := strconv.Atoi(heightStr)

		if errW != nil || errH != nil || width <= 0 || height <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Width dan Height harus berupa angka positif",
			})
			return
		}
		if width > h.config.Image.MaxDimension || height > h.config.Image.MaxDimension {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Width dan Height maksimal %d px", h.config.Image.MaxDimension),
			})
			return
		}

		safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(file.Filename))
		srcPath := filepath.Join(h.config.Storage.UploadDir, safeFilename)

		os.MkdirAll(h.config.Storage.UploadDir, 0755)

		if err := c.SaveUploadedFile(file, srcPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Gagal menyimpan file", "error": err.Error()})
			return
		}
		defer os.Remove(srcPath)

		resultPath, err := services.ProcessImageResize(srcPath, width, height)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Gagal resize gambar", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Resize berhasil",
			"data": gin.H{
				"original_file": file.Filename,
				"result_path":   resultPath,
				"width":         width,
				"height":        height,
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
