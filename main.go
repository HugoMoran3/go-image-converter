package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/bmp"
	"github.com/mat/besticon/ico"
	"github.com/tdewolff/canvas"
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

}
