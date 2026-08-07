package manager

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// saveUploadedImage saves an uploaded file under uploadsDir/subdir/slotName.ext,
// deleting any prior file with the same slot name (regardless of extension)
// first — so replacing a logo doesn't leave the old file orphaned on disk.
func saveUploadedImage(uploadsDir string, subdir string, slotName string, fileHeader *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	targetDir := filepath.Join(uploadsDir, subdir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	matches, _ := filepath.Glob(filepath.Join(targetDir, slotName+".*"))
	for _, m := range matches {
		_ = os.Remove(m)
	}

	destPath := filepath.Join(targetDir, slotName+ext)
	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("create destination file: %w", err)
	}
	defer dest.Close()

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return "", fmt.Errorf("copy uploaded file: %w", err)
	}

	return "/uploads/" + subdir + "/" + slotName + ext, nil
}
