package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".svg": true, ".avif": true,
}

var (
	ErrFileTooLarge = errors.New("file exceeds maximum size")
	ErrInvalidImage = errors.New("invalid image format")
	ErrInvalidFile  = errors.New("invalid file")
)

func ValidateImageFile(header *multipart.FileHeader, maxSize int64) (string, error) {
	if header == nil {
		return "", ErrInvalidImage
	}
	if header.Size > maxSize {
		return "", ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExts[ext] {
		return "", ErrInvalidImage
	}

	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	mime := http.DetectContentType(buf[:n])

	if ext == ".svg" {
		if !strings.HasPrefix(mime, "text/") && mime != "application/xml" {
			return "", ErrInvalidImage
		}
	} else if !strings.HasPrefix(mime, "image/") {
		return "", ErrInvalidImage
	}

	return ext, nil
}

func SaveUploadedImage(header *multipart.FileHeader, dir string, maxSize int64) (string, error) {
	ext, err := ValidateImageFile(header, maxSize)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s-%d%s", uuid.NewString(), time.Now().UnixMilli(), ext)
	dst := filepath.Join(dir, filename)

	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}

	return dst, nil
}

func SaveUploadedFile(header *multipart.FileHeader, dir string, maxSize int64) (string, error) {
	if header == nil {
		return "", ErrInvalidFile
	}
	if header.Size > maxSize {
		return "", ErrFileTooLarge
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := fmt.Sprintf("%s-%d%s", uuid.NewString(), time.Now().UnixMilli(), ext)
	dst := filepath.Join(dir, filename)

	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}

	return filename, nil
}
