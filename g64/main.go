package main

import (
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fyne-io/image/ico"
	"github.com/jackmordaunt/icns"
	"golang.design/x/clipboard"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go64 <image_file>")
		os.Exit(1)
	}

	imagePath := os.Args[1]

	err := clipboard.Init()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	dataURI, err := convertToDataURI(imagePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}
	clipboard.Write(clipboard.FmtText, []byte(dataURI))
}

func convertToDataURI(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	mimeType, err := determineMimeType(imagePath, imageData)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(imageData)
	dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
	return dataURI, nil
}

func determineMimeType(imagePath string, imageData []byte) (string, error) {
	ext := strings.ToLower(filepath.Ext(imagePath))
	switch ext {
	case ".png":
		_, err := png.Decode(strings.NewReader(string(imageData)))
		if err == nil {
			return "image/png", nil
		}
	case ".jpg", ".jpeg":
		_, err := jpeg.Decode(strings.NewReader(string(imageData)))
		if err == nil {
			return "image/jpeg", nil
		}
	case ".ico":
		_, err := ico.Decode(strings.NewReader(string(imageData)))
		if err == nil {
			return "image/x-icon", nil
		}
	case ".icns":
		_, err := icns.Decode(strings.NewReader(string(imageData)))
		if err == nil {
			return "image/icns", nil
		}
	}
	return "", fmt.Errorf("unsupported file extension: %s", ext)
}
