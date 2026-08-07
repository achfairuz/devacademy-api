package utils

import (
	"image/jpeg"

	"os"
)

func Compress(inputFile string) (string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return "", err
	}
	OutputFile, err := os.Create(inputFile + ".webp")
	if err != nil {
		return "", err
	}
	defer OutputFile.Close()

	err = jpeg.Encode(OutputFile, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return "", err
	}
	return OutputFile.Name(), nil
}
