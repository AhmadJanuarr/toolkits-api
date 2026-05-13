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

func (h *ImageHandler) ConvertImage(w http.ResponseWriter, r *http.Request) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		targetFormat := strings.ToLower(r.FormValue("format"))
		if targetFormat == "" || !utils.Contains(h.config.Image.AllowedFormats, targetFormat) {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Format target tidak valid atau tidak didukung (gunakan: jpg, jpeg, png, webp)",
			})
			return
		}

		_, fileHeader, err := r.FormFile("file")
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "File tidak ditemukan atau tidak valid",
				"error":   err.Error(),
			})
			return
		}

		if err := utils.ValidateImageFile(fileHeader); err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		if fileHeader.Size > h.config.Image.MaxFileSize {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Ukuran file tidak boleh lebih dari %d MB", h.config.Image.MaxFileSize>>20),
			})
			return
		}

		safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))
		srcPath := filepath.Join(h.config.Storage.TempDir, safeFilename)

		if err := utils.SaveUploadedFile(fileHeader, srcPath); err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Gagal menyimpan file sementara",
				"error":   err.Error(),
			})
			return
		}
		defer os.Remove(srcPath)
		resultPath, err := services.ProcessImageConversion(srcPath, targetFormat)
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Gagal mengonversi gambar",
				"error":   err.Error(),
			})
			return
		}

		utils.DownloadFile(w, r, resultPath)

	default:
		utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}
}

func (h *ImageHandler) CompressionImage(w http.ResponseWriter, r *http.Request) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		_, fileHeader, err := r.FormFile("file")
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "File tidak ditemukan atau tidak valid",
				"error":   err.Error(),
			})
			return
		}

		qualityStr := r.FormValue("quality")
		if qualityStr == "" {
			qualityStr = "80"
		}

		quality, err := strconv.Atoi(qualityStr)
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Quality harus berupa angka",
				"error":   err.Error(),
			})
			return
		}
		if quality < h.config.Image.MinQuality || quality > h.config.Image.MaxQuality {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Quality harus antara %d dan %d", h.config.Image.MinQuality, h.config.Image.MaxQuality),
			})
			return
		}

		safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))
		srcPath := filepath.Join(h.config.Storage.TempDir, safeFilename)

		if err := utils.SaveUploadedFile(fileHeader, srcPath); err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Gagal menyimpan file sementara",
				"error":   err.Error(),
			})
			return
		}
		defer os.Remove(srcPath)

		resultPath, err := services.ProcessImageCompression(srcPath, quality)
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Gagal mengompres gambar",
				"error":   err.Error(),
			})
			return
		}

		utils.DownloadFile(w, r, resultPath)

	default:
		utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}
}

func (h *ImageHandler) ResizeImage(w http.ResponseWriter, r *http.Request) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		_, fileHeader, err := r.FormFile("file")
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "File wajib diupload",
				"error":   err.Error(),
			})
			return
		}

		widthStr := r.FormValue("width")
		heightStr := r.FormValue("height")

		width, errW := strconv.Atoi(widthStr)
		height, errH := strconv.Atoi(heightStr)

		if errW != nil || errH != nil || width <= 0 || height <= 0 {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Width dan Height harus berupa angka positif",
			})
			return
		}
		if width > h.config.Image.MaxDimension || height > h.config.Image.MaxDimension {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": fmt.Sprintf("Width dan Height maksimal %d px", h.config.Image.MaxDimension),
			})
			return
		}

		safeFilename := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))
		srcPath := filepath.Join(h.config.Storage.UploadDir, safeFilename)

		os.MkdirAll(h.config.Storage.UploadDir, 0755)

		if err := utils.SaveUploadedFile(fileHeader, srcPath); err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Gagal menyimpan file",
				"error":   err.Error(),
			})
			return
		}
		defer os.Remove(srcPath)

		resultPath, err := services.ProcessImageResize(srcPath, width, height)
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Gagal resize gambar",
				"error":   err.Error(),
			})
			return
		}

		utils.DownloadFile(w, r, resultPath)

	default:
		utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}
}
