package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg" // Register JPEG decoder
	"image/png"
	_ "image/png" // Register PNG decoder
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Supported image extensions (case-insensitive)
var supportedExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

type imageInfo struct {
	Path string
	Img  image.Image
}

func main() {
	// --- Command Line Flags ---
	inputDir := flag.String("d", "", "Directory containing input images (required)")
	outputFile := flag.String("o", "stitched_output.png", "Output image file path (defaults to stitched_output.png)")
	cols := flag.Int("cols", 4, "Number of columns in the grid layout")
	// Print Help Message
	flag.Usage = func() {
		fmt.Println("Usage: stitch -d <directory> -o <output.whatever> -cols <integer>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *inputDir == "" {
		log.Fatal("Error: Input directory (--dir) is required.")
		flag.Usage()
		os.Exit(1)
	}
	if *cols <= 0 {
		log.Fatal("Error: Number of columns (--cols) must be positive.")
	}

	// --- Find and Load Images ---
	images, err := loadImagesFromDir(*inputDir)
	if err != nil {
		log.Fatalf("Error loading images: %v", err)
	}
	if len(images) == 0 {
		log.Fatalf("No supported image files found in directory: %s", *inputDir)
	}
	log.Printf("Found %d images to stitch.", len(images))

	// --- Calculate Layout and Canvas Size ---
	totalWidth, totalHeight, imagePositions := calculateLayout(images, *cols)
	log.Printf("Calculated canvas size: %d x %d pixels", totalWidth, totalHeight)

	// --- Create Canvas and Draw Images ---
	// Create a new RGBA image. RGBA allows for transparency if needed,
	// although here we are just copying opaque images.
	canvas := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Optional: Fill background if you don't want default black/transparent
	// bg := image.Uniform{color.White}
	// draw.Draw(canvas, canvas.Bounds(), &bg, image.Point{}, draw.Src)

	log.Println("Drawing images onto canvas...")
	for i, imgInfo := range images {
		pos := imagePositions[i]
		bounds := imgInfo.Img.Bounds()
		// Define the rectangle on the canvas where this image will be drawn
		destRect := image.Rect(pos.X, pos.Y, pos.X+bounds.Dx(), pos.Y+bounds.Dy())

		// Draw the image onto the canvas
		// draw.Src copies pixel data directly.
		// draw.Over handles transparency (useful if input images have alpha)
		draw.Draw(canvas, destRect, imgInfo.Img, bounds.Min, draw.Over)
		log.Printf("  Drew %s at (%d, %d)", filepath.Base(imgInfo.Path), pos.X, pos.Y)
	}

	// --- Save the Output Image ---
	log.Printf("Saving stitched image to %s...", *outputFile)
	outFile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Error creating output file %s: %v", *outputFile, err)
	}
	defer outFile.Close()

	// Encode based on output file extension
	outputExt := strings.ToLower(filepath.Ext(*outputFile))
	var encodeErr error // Use a local variable for the encoding error

	switch outputExt {
	case ".jpg", ".jpeg":
		log.Println("Encoding as JPEG (Quality 100)...")
		// Ensure the jpeg package was imported: import _ "image/jpeg"
		// jpeg.Encode expects: io.Writer, image.Image, *jpeg.Options
		// Using high quality; 90-95 is often a good balance vs file size.
		options := &jpeg.Options{Quality: 100}
		encodeErr = jpeg.Encode(outFile, canvas, options)

	case ".png":
		log.Println("Encoding as PNG...")
		// Ensure the png package was imported: import _ "image/png"
		// png.Encode expects: io.Writer, image.Image
		encodeErr = png.Encode(outFile, canvas)

	default:
		// If the extension is not recognized, default to PNG
		log.Printf("Warning: Unsupported output extension '%s'. Saving as PNG.", outputExt)
		log.Println("Encoding as PNG...")
		// Ensure the png package was imported: import _ "image/png"
		encodeErr = png.Encode(outFile, canvas)
	}

	// Check for encoding errors after the switch statement
	if encodeErr != nil {
		// It's often good practice to attempt removing the partially written/failed file
		outFile.Close()        // Close before removing
		os.Remove(*outputFile) // Attempt removal, ignore error if it fails
		log.Fatalf("Error encoding output image to %s: %v", *outputFile, encodeErr)
	}

	log.Println("Successfully created stitched image:", *outputFile)
} // Assuming this brace closes the main function

// loadImagesFromDir scans a directory, decodes supported images, and returns them.
func loadImagesFromDir(dirPath string) ([]imageInfo, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dirPath, err)
	}

	var images []imageInfo
	log.Printf("Scanning directory: %s", dirPath)

	for _, file := range files {
		if file.IsDir() {
			continue // Skip subdirectories
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if !supportedExtensions[ext] {
			continue // Skip unsupported files
		}

		fullPath := filepath.Join(dirPath, file.Name())
		log.Printf("  Loading: %s", file.Name())

		imgFile, err := os.Open(fullPath)
		if err != nil {
			log.Printf("    Warning: Could not open file %s: %v", file.Name(), err)
			continue // Skip files that can't be opened
		}
		// No defer here, Decode might keep the file open, close explicitly below

		img, format, err := image.Decode(imgFile)
		imgFile.Close() // Close the file handle *after* decoding
		if err != nil {
			log.Printf("    Warning: Could not decode image %s (format: %s): %v", file.Name(), format, err)
			continue // Skip files that can't be decoded
		}

		log.Printf("    Decoded %s (format: %s, size: %dx%d)", file.Name(), format, img.Bounds().Dx(), img.Bounds().Dy())
		images = append(images, imageInfo{Path: fullPath, Img: img})
	}

	return images, nil
}

// calculateLayout determines the grid positions and total canvas size.
func calculateLayout(images []imageInfo, cols int) (totalWidth, totalHeight int, positions []image.Point) {
	if len(images) == 0 || cols <= 0 {
		return 0, 0, nil
	}

	positions = make([]image.Point, len(images))
	currentX := 0
	currentY := 0
	rowMaxHeight := 0
	maxSeenWidth := 0 // Track the maximum width reached in any row

	for i, imgInfo := range images {
		bounds := imgInfo.Img.Bounds()
		imgWidth := bounds.Dx()
		imgHeight := bounds.Dy()

		// Calculate position for the current image
		positions[i] = image.Point{X: currentX, Y: currentY}

		// Update current X position for the next image in the row
		currentX += imgWidth

		// Update the maximum height encountered in the current row
		if imgHeight > rowMaxHeight {
			rowMaxHeight = imgHeight
		}

		// Update the overall maximum width seen so far
		if currentX > maxSeenWidth {
			maxSeenWidth = currentX
		}

		// Check if we need to move to the next row
		// (i+1) because 'i' is 0-based index
		if (i+1)%cols == 0 {
			// Move to the next row
			currentY += rowMaxHeight // Add the max height of the completed row
			currentX = 0             // Reset X for the new row
			rowMaxHeight = 0         // Reset max height for the new row
		}
	}

	// After the loop, if the last row wasn't full, we need to add its height
	if len(images)%cols != 0 {
		totalHeight = currentY + rowMaxHeight
	} else {
		// If the last row was full, currentY already includes the last row's max height
		totalHeight = currentY
	}

	totalWidth = maxSeenWidth // The total width is the widest row encountered

	return totalWidth, totalHeight, positions
}
