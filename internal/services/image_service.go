package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"toolkits/internal/utils"

	"github.com/disintegration/imaging"
)

func ProcessImageConversion(inputPath string, targetFormat string) (string, error) {

	// input
	img, _, err := utils.LoadImage(inputPath)
	if err != nil {
		return "", err
	}

	// process
	filename := filepath.Base(inputPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// output
	outputDir := "./temp/processed"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}
	newFilename := fmt.Sprintf("%s.%s", nameWithoutExt, strings.ToLower(targetFormat))
	outputPath := filepath.Join(outputDir, newFilename)

	if err := utils.SaveImage(img, outputPath, targetFormat, 80); err != nil {
		return "", err
	}

	return outputPath, nil
}

func ProcessImageCompression(inputPath string, quality int) (string, error) {

	img, format, err := utils.LoadImage(inputPath)
	if err != nil {
		return "", err
	}

	outputDir := "./temp/compressed"
	os.MkdirAll(outputDir, 0755)

	filename := filepath.Base(inputPath)
	newFilename := fmt.Sprintf("%s_compressed.%s", strings.TrimSuffix(filename, filepath.Ext(filename)), format)
	outputPath := filepath.Join(outputDir, newFilename)

	if err := utils.SaveImage(img, outputPath, format, quality); err != nil {
		return "", err
	}
	return outputPath, nil
}

func ProcessImageResize(inputPath string, width int, height int) (string, error) {

	img, format, err := utils.LoadImage(inputPath)
	if err != nil {
		return "", err
	}
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	outputDir := "./temp/resized"
	os.MkdirAll(outputDir, 0755)

	filename := filepath.Base(inputPath)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	newFilename := fmt.Sprintf("%s_resized_%dx%d.%s", nameWithoutExt, width, height, format)
	outputPath := filepath.Join(outputDir, newFilename)

	if err := utils.SaveImage(resizedImg, outputPath, format, 90); err != nil {
		return "", err
	}
	return outputPath, nil
}
