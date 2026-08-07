package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/deepteams/webp"
)

func Compress(inputFile string) (string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return inputFile, nil
	}

	if format == "webp" {
		return inputFile, nil
	}

	outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "-compressed.webp"

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, webp.DefaultOptions()); err != nil {
		return "", err
	}

	if err := os.WriteFile(outputFile, buf.Bytes(), 0o644); err != nil {
		return "", err
	}

	if err := os.Remove(inputFile); err != nil {
		return "", err
	}

	return outputFile, nil
}
