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

// saveUploadedImage validates and saves an uploaded image file under
// uploadsDir/subdir, using slotName as the fixed base filename (so a
// re-upload always overwrites the same slot). Any existing file(s) for
// that slot — regardless of prior extension — are deleted first, so
// stale images never accumulate on disk.
// Returns the public URL path the file is served at (e.g. "/uploads/homepage/image_one.png").
func saveUploadedImage(uploadsDir string, subdir string, slotName string, fileHeader *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("unsupported file type %q (allowed: jpg, jpeg, png, webp)", ext)
	}

	dir := filepath.Join(uploadsDir, subdir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	// Delete any existing file(s) for this slot, regardless of extension,
	// so replacing a .png with a .jpg doesn't leave the old .png orphaned.
	existing, _ := filepath.Glob(filepath.Join(dir, slotName+".*"))
	for _, old := range existing {
		_ = os.Remove(old)
	}

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
