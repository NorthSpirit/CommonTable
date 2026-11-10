package admin

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func validateImageFile(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedImageExtensions[ext] {
		return fmt.Errorf("unsupported file type %q (allowed: jpg, jpeg, png, webp)", ext)
	}
	return nil
}

func saveUploadedImage(uploadsDir string, subdir string, slotName string, fileHeader *multipart.FileHeader) (string, error) {
	if err := validateImageFile(fileHeader); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	dir := filepath.Join(uploadsDir, subdir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	removeUploadedImages(uploadsDir, subdir, slotName)

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	destPath := filepath.Join(dir, slotName+ext)
	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("create destination file: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return "", fmt.Errorf("write uploaded file: %w", err)
	}

	return "/uploads/" + subdir + "/" + slotName + ext, nil
}

func removeUploadedImages(uploadsDir string, subdir string, slotName string) {
	existing, _ := filepath.Glob(filepath.Join(uploadsDir, subdir, slotName+".*"))
	for _, old := range existing {
		_ = os.Remove(old)
	}
}
