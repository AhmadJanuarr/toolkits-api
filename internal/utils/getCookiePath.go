package utils

import (
	"os"
	"path/filepath"
)

func GetCookiePath(platform string) string {
	var sourcePath, cookieName string

	if platform == "youtube" {
		sourcePath = "/etc/secrets/youtube_cookies.txt"
		cookieName = "youtube_cookies.txt"
	} else if platform == "instagram" {
		sourcePath = "/etc/secrets/ig-cookies.txt"
		cookieName = "ig-cookies.txt"
	} else {
		return ""
	}

	// Cek apakah file asli ada di /etc/secrets
	if _, err := os.Stat(sourcePath); err != nil {
		// Jika tidak ada (misal di environment lokal), kembalikan path aslinya saja
		return sourcePath
	}

	// Buat path tujuan di folder temporary sistem (biasanya /tmp)
	destPath := filepath.Join(os.TempDir(), cookieName)

	// Baca dari source yang read-only dan tulis ke destPath yang writable
	data, err := os.ReadFile(sourcePath)
	if err == nil {
		_ = os.WriteFile(destPath, data, 0644)
	}

	return destPath
}
