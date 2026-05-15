package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"toolkits/internal/config"
	"toolkits/internal/services"
	"toolkits/internal/utils"
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

func (h *DownloaderHandler) Downloader(w http.ResponseWriter, r *http.Request) {
	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()

		inputURL := r.FormValue("inputURL")
		formatID := r.FormValue("format_id")
		if inputURL == "" || formatID == "" {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Input dan FormatID tidak boleh kosong",
			})
			return
		}
		cacheKey := inputURL + "_" + formatID
		if cachedPath, exists := utils.FileCache.Get(cacheKey); exists {
			fileName := filepath.Base(cachedPath)
			w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
			http.ServeFile(w, r, cachedPath)
			return
		}
		resultPath, err := services.ProsessDownload(r.Context(), inputURL, formatID)
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Maaf, download gagal dikarenakan server waktu tunggu melebihi dari 30 detik, silahkan coba lagi",
				"error":   err.Error(),
			})
			return
		}
		utils.FileCache.Set(cacheKey, resultPath)
		http.SetCookie(w, &http.Cookie{
			Name:  "download_ready",
			Value: "1",
			Path:  "/",
		})
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(resultPath)+"\"")
		http.ServeFile(w, r, resultPath)
	default:
		utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
	}
}

func (h *DownloaderHandler) DownloaderGetInfo(w http.ResponseWriter, r *http.Request) {

	select {
	case h.semaphore <- struct{}{}:
		defer func() { <-h.semaphore }()
		//
		inputURL := r.FormValue("inputURL")
		if strings.Contains(inputURL, "/photo/") {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "Maaf, untuk saat ini format tiktok berupa slide foto belum mendukung",
			})
			return
		}

		if inputURL == "" {
			utils.JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"status":  http.StatusBadRequest,
				"message": "inputURL tidak boleh kosong atau tidak valid",
			})
			return
		}

		info, err := services.ProcessGetInfo(r.Context(), inputURL)
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": "Gagal mengambil informasi video, waktu tunggu melebihi batas 30 detik, silahkan coba lagi",
				"error":   err.Error(),
			})
			return
		}
		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "Informasi video berhasil diambil",
			"data": map[string]interface{}{
				"info": info,
			},
		})
	default:
		utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  http.StatusServiceUnavailable,
			"message": "Server sibuk, coba lagi sebentar",
		})
		return
	}

}
