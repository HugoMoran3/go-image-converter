package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hhrutter/tiff"
	"golang.org/x/image/webp"
)

var SupportedFormats = map[string]struct{}{
	"png": {},
	"jpg": {},
	"jpeg": {},
	"webp": {},
	"tiff": {},
	"bmp": {},
	"ico": {},
}

// convertFromSVG handles SVG input. Currently a stub.
func convertFromSVG(inputPath, outputPath string, quality, icoSize int) error {
	// TODO: Implement SVG conversion
	// For now, let's return an error indicating it's not implemented if an output path is given
	if outputPath != "" {
		return fmt.Errorf("SVG conversion to %s not implemented yet", filepath.Ext(outputPath))
	}
	log.Println("SVG processing stub called for input:", inputPath) // Or handle as pure input
	return nil
}

// convertImageLogic contains the core image conversion functionality.
func convertImageLogic(inputPath, outputPath string, quality, icoSize int) error {
	// SVG processing (moved here from main)
	if strings.ToLower(filepath.Ext(inputPath)) == ".svg" {
		return convertFromSVG(inputPath, outputPath, quality, icoSize)
	}

	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %w", inputPath, err)
	}
	defer inputFile.Close()

	// Decode image based on input type
	var img image.Image
	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".webp":
		img, err = webp.Decode(inputFile)
		if err != nil {
			return fmt.Errorf("error decoding WebP image %s: %w", inputPath, err)
		}
	case ".tiff", ".tif":
		img, err = tiff.Decode(inputFile)
		if err != nil {
			return fmt.Errorf("error decoding TIFF image %s: %w", inputPath, err)
		}
	default:
		img, _, err = image.Decode(inputFile) // The '_' was for format name, not strictly needed here
		if err != nil {
			return fmt.Errorf("error decoding image %s: %w", inputPath, err)
		}
	}

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	// Encode the image to the output file
	switch strings.ToLower(filepath.Ext(outputPath)) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: quality})
	case ".png":
		err = png.Encode(outputFile, img)
	case ".tiff", ".tif":
		err = tiff.Encode(outputFile, img, &tiff.Options{
			Compression: tiff.LZW, // Example option, adjust as needed
		})
	// Note: WEBP encoding is not included in golang.org/x/image/webp by default.
	// If WEBP output is desired, a different library or cgo bindings to libwebp would be needed.
	// case ".webp":
	// 	 return fmt.Errorf("WEBP output encoding not supported by standard library package")
	default:
		return fmt.Errorf("unsupported output format: %s", filepath.Ext(outputPath))
	}

	if err != nil {
		return fmt.Errorf("error encoding to %s: %w", outputPath, err)
	}

	return nil
}

func main() {
	// Define command line flags
	inputPath := flag.String("input", "", "Path to input file")
	outputPath := flag.String("output", "", "Path to output file (format will be inferred from extension)")
	quality := flag.Int("quality", 90, "JPEG compression quality (0-100)")
	icoSize := flag.Int("ico", 32, "Icon size for ICO format (16, 32, 48 or 68)")

	flag.Parse()

	// Validate input parameters
	if *inputPath == "" || *outputPath == "" { // Also ensure output path is provided for conversion
		fmt.Fprintf(os.Stderr, "Usage: %s -input <inputFile> -output <outputFile> [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := convertImageLogic(*inputPath, *outputPath, *quality, *icoSize)
	if err != nil {
		log.Printf("Conversion failed: %v", err) // Use log.Printf to print to stderr
		os.Exit(1)
	}

	log.Println("Conversion successful.")
}