package utils

import (
	"net/http"
	"path/filepath"
)

func DownloadFile(w http.ResponseWriter, r *http.Request, filePath string) {
	fileName := filepath.Base(filePath)
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fileName+"\"")
	http.ServeFile(w, r, filePath)
}
