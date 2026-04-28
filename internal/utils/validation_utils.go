package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
)

// ValidateImageFile mengecek apakah file benar-benar gambar valid menggunakan Magic Bytes
func ValidateImageFile(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file upload: %v", err)
	}
	defer src.Close()

	// Baca 512 byte pertama untuk sniffing
	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return fmt.Errorf("gagal membaca header file: %v", err)
	}

	// PENTING: Reset pointer file ke awal setelah dibaca
	if _, err := src.Seek(0, 0); err != nil {
		return fmt.Errorf("gagal reset pointer file: %v", err)
	}

	// Deteksi tipe konten
	contentType := http.DetectContentType(buffer)

	// Whitelist tipe yang diizinkan
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	if !validTypes[contentType] {
		return fmt.Errorf("tipe file tidak valid atau berbahaya: %s (hanya mendukung JPG, PNG, WEBP)", contentType)
	}

	return nil
}
