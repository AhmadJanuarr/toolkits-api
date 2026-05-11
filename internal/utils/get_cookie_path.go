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

	if _, err := os.Stat(sourcePath); err != nil {
		if _, errLocal := os.Stat(cookieName); errLocal == nil {
			absPath, _ := filepath.Abs(cookieName)
			return absPath
		}
		return ""
	}

	destPath := filepath.Join(os.TempDir(), cookieName)
	data, err := os.ReadFile(sourcePath)
	if err == nil {
		_ = os.WriteFile(destPath, data, 0644)
	}

	return destPath
}
