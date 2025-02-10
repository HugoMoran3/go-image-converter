package main

import (
	"flag"
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
// Add the missing SVG conversion function
func convertFromSVG(inputPath, outputPath string, quality, icoSize int) error {
	// TODO: Implement SVG conversion
	return nil
}

func main() {
	// Define command line flags
	inputPath := flag.String("input", "", "Path to input file")
	outputPath := flag.String("output", "", "Path to output file (format will be inferred from extension)")
	quality := flag.Int("quality", 90, "JPEG compression quality (0-100)")
	icoSize := flag.Int("ico", 32, "Icon size for ICO format (16, 32, 48 or 68)")

	flag.Parse()

	// Validate input perameters
	if *inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// SVG processing
	if strings.ToLower(filepath.Ext(*inputPath)) == ".svg" {
		convertFromSVG(*inputPath, *outputPath, *quality, *icoSize)
		return
	}

	// Open input file
	inputFile, err := os.Open(*inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// WebP processing
	var img image.Image
	if ext := strings.ToLower(filepath.Ext(*inputPath)); ext == ".webp" {
		img, err = webp.Decode(inputFile)
		if err != nil {
			log.Fatalf("Error decoding WebP image: %v", err)
		}
	} else if ext == ".tiff" || ext == ".tif" {
		img, err = tiff.Decode(inputFile)
		if err != nil {
			log.Fatalf("Error decoding TIFF image: %v", err)
		}
	} else {
		img, _, err = image.Decode(inputFile)
		if err != nil {
			log.Fatalf("Error decoding image: %v", err)
		}
	}

	// Create the output file
	outputFile, err := os.Create(*outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Add this: Encode the image to the output file
	switch strings.ToLower(filepath.Ext(*outputPath)) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: *quality})
	case ".png":
		err = png.Encode(outputFile, img)
	case ".tiff", ".tif":
		err = tiff.Encode(outputFile, img, &tiff.Options{
			Compression: tiff.LZW,
		})
	default:
		log.Fatal("Unsupported output format")
	}
	if err != nil {
		log.Fatal(err)
	}
}