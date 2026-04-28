package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"toolkits/internal/models"
)

type ytdlpJSON struct {
	Title     string  `json:"title"`
	Uploader  string  `json:"uploader"`
	Duration  float64 `json:"duration"`
	Thumbnail string  `json:"thumbnail"`
	Formats   []struct {
		FormatID       string   `json:"format_id"`
		Ext            string   `json:"ext"`
		FormatNote     string   `json:"format_note"`
		Vcodec         string   `json:"vcodec"`
		Acodec         string   `json:"acodec"`
		Filesize       *float64 `json:"filesize"`
		FilesizeApprox *float64 `json:"filesize_approx"`
	} `json:"formats"`
}

func getPlatformFromURL(inputURL string) string {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "general"
	}
	hostname := strings.ToLower(parsedURL.Host)

	if strings.Contains(hostname, "youtube.com") || strings.Contains(hostname, "youtu.be") {
		return "youtube"
	} else if strings.Contains(hostname, "tiktok.com") {
		return "tiktok"
	} else if strings.Contains(hostname, "instagram.com") {
		return "instagram"
	} else {
		return "general"
	}
}

func ProcessGetInfo(inputURL string) (*models.InfoResponse, error) {

	platform := getPlatformFromURL(inputURL)
	args := []string{"-J", "--no-playlist", "--no-warnings", "--force-ipv4", inputURL}

	if platform == "tiktok" && strings.Contains(inputURL, "/photo/") {
		return nil, fmt.Errorf("Maaf, untuk saat ini format tiktok berupa slide foto belum didukung")
	}

	if platform == "youtube" {
		args = append([]string{"--cookies", "youtube_cookies.txt", "--js-runtimes", "node"}, args...)
	} else if platform == "instagram" {
		args = append([]string{"--cookies", "ig-cookies.txt"}, args...)
	}

	cmd := exec.Command("yt-dlp", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gagal ambil info video: %v, log: %s", err, stderr.String())
	}

	var raw ytdlpJSON
	if err := json.Unmarshal(out.Bytes(), &raw); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON metadata: %v", err)
	}

	info := &models.InfoResponse{
		Title:     raw.Title,
		Author:    raw.Uploader,
		Duration:  fmt.Sprintf("%d detik", int(raw.Duration)),
		Thumbnail: raw.Thumbnail,
	}

	for _, format := range raw.Formats {
		if format.Vcodec == "none" && format.Acodec == "none" {
			continue
		}
		hasVideo := format.Vcodec != "none"
		hasAudio := format.Acodec != "none"
		quality := format.FormatNote
		var finalFilesize int64 = 0

		if format.Filesize != nil {
			finalFilesize = int64(*format.Filesize)
		} else if format.FilesizeApprox != nil {
			finalFilesize = int64(*format.FilesizeApprox)
		}

		if quality == "" && platform == "tiktok" {
			quality = "original"
		} else if !hasVideo && hasAudio {
			quality = "Audio only"
		}

		info.Formats = append(info.Formats, models.FormatOption{
			FormatID: format.FormatID,
			Quality:  quality,
			MimeType: format.Ext,
			HasVideo: hasVideo,
			HasAudio: hasAudio,
			Filesize: &finalFilesize,
		})
	}

	return info, nil
}

func ProsessDownload(inputURL string, formatID string) (string, error) {

	platform := getPlatformFromURL(inputURL)
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan direktori saat ini: %v", err)
	}
	outputDir := filepath.Join(currentDir, "temp", "downloads", platform)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	templateName := "%(title)s_%(id)s.%(ext)s"

	args := []string{
		"-f", formatID,
		"-P", outputDir,
		"-o", templateName,
		"--print", "filename",
		"--no-simulate",
		"--no-mtime",
		"--no-warnings",
		"--force-ipv4",
		inputURL}

	if platform == "youtube" {
		args = append([]string{"--cookies", "youtube_cookies.txt", "--js-runtimes", "node"}, args...)
	} else if platform == "instagram" {
		args = append([]string{"--cookies", "ig-cookies.txt"}, args...)
	}

	cmd := exec.Command("yt-dlp", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("gagal mengunduh video: %v, log: %s", err, stderr.String())
	}

	finalPath := strings.TrimSpace(out.String())
	return finalPath, nil
}
