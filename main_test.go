package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// TestMain is used for setup and teardown of tests
func TestMain(m *testing.M) {
	// Setup
	err := os.MkdirAll("test_files/output", 0755)
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll("test_files/output")
	
	os.Exit(code)
}

// Basic table-driven test
func TestImageConversion(t *testing.T) {
	// Save original args and flags
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	tests := []struct {
		name        string
		inputPath   string
		outputPath  string
		quality     int
		expectError bool
	}{
		// JPEG conversions
		{name: "JPG to PNG", inputPath: "input/jpg_example.jpg", outputPath: "test_files/output/jpg_to_png.png", quality: 90},
		{name: "JPG to TIFF", inputPath: "input/jpg_example.jpg", outputPath: "test_files/output/jpg_to_tiff.tiff", quality: 90},
		{name: "JPG to WEBP", inputPath: "input/jpg_example.jpg", outputPath: "test_files/output/jpg_to_webp.webp", quality: 90},

		// PNG conversions
		{name: "PNG to JPG", inputPath: "input/png_example.png", outputPath: "test_files/output/png_to_jpg.jpg", quality: 90},
		{name: "PNG to TIFF", inputPath: "input/png_example.png", outputPath: "test_files/output/png_to_tiff.tiff", quality: 90},
		{name: "PNG to WEBP", inputPath: "input/png_example.png", outputPath: "test_files/output/png_to_webp.webp", quality: 90},

		// WEBP conversions
		{name: "WEBP to JPG", inputPath: "input/webp_exmaple.webp", outputPath: "test_files/output/webp_to_jpg.jpg", quality: 90},
		{name: "WEBP to PNG", inputPath: "input/webp_exmaple.webp", outputPath: "test_files/output/webp_to_png.png", quality: 90},
		{name: "WEBP to TIFF", inputPath: "input/webp_exmaple.webp", outputPath: "test_files/output/webp_to_tiff.tiff", quality: 90},

		// TIFF conversions
		{name: "TIFF to JPG", inputPath: "input/tiff_example.tiff", outputPath: "test_files/output/tiff_to_jpg.jpg", quality: 90},
		{name: "TIFF to PNG", inputPath: "input/tiff_example.tiff", outputPath: "test_files/output/tiff_to_png.png", quality: 90},
		{name: "TIFF to WEBP", inputPath: "input/tiff_example.tiff", outputPath: "test_files/output/tiff_to_webp.webp", quality: 90},

		// BMP conversions
		{name: "BMP to JPG", inputPath: "input/bmp_example.bmp", outputPath: "test_files/output/bmp_to_jpg.jpg", quality: 90},
		{name: "BMP to PNG", inputPath: "input/bmp_example.bmp", outputPath: "test_files/output/bmp_to_png.png", quality: 90},
		{name: "BMP to WEBP", inputPath: "input/bmp_example.bmp", outputPath: "test_files/output/bmp_to_webp.webp", quality: 90},
		{name: "BMP to TIFF", inputPath: "input/bmp_example.bmp", outputPath: "test_files/output/bmp_to_tiff.tiff", quality: 90},

		// SVG conversions (generally expect errors for raster output unless specifically handled)
		{name: "SVG to PNG", inputPath: "input/svg_example.svg", outputPath: "test_files/output/svg_to_png.png", quality: 90, expectError: true},
		{name: "SVG to JPG", inputPath: "input/svg_example.svg", outputPath: "test_files/output/svg_to_jpg.jpg", quality: 90, expectError: true},
		{name: "SVG to TIFF", inputPath: "input/svg_example.svg", outputPath: "test_files/output/svg_to_tiff.tiff", quality: 90, expectError: true},
		{name: "SVG to WEBP", inputPath: "input/svg_example.svg", outputPath: "test_files/output/svg_to_webp.webp", quality: 90, expectError: true},

		// ICO conversions (generally expect errors for conversion to/from other formats unless specifically handled)
		{name: "ICO to PNG", inputPath: "input/favicon_example.ico", outputPath: "test_files/output/ico_to_png.png", quality: 90, expectError: true},
		{name: "ICO to JPG", inputPath: "input/favicon_example.ico", outputPath: "test_files/output/ico_to_jpg.jpg", quality: 90, expectError: true},
		{name: "ICO to TIFF", inputPath: "input/favicon_example.ico", outputPath: "test_files/output/ico_to_tiff.tiff", quality: 90, expectError: true},
		{name: "ICO to WEBP", inputPath: "input/favicon_example.ico", outputPath: "test_files/output/ico_to_webp.webp", quality: 90, expectError: true},

		// Cross-conversions from other types to ICO (expectError: true)
		{name: "PNG to ICO", inputPath: "input/png_example.png", outputPath: "test_files/output/png_to_ico.ico", quality: 90, expectError: true},
		{name: "JPG to ICO", inputPath: "input/jpg_example.jpg", outputPath: "test_files/output/jpg_to_ico.ico", quality: 90, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			
			// Reset output path for each test
			os.RemoveAll(tt.outputPath)

			// Run the conversion using your main package functions
			args := []string{
				"-input", tt.inputPath,
				"-output", tt.outputPath,
				"-quality", fmt.Sprintf("%d", tt.quality),
			}
			
			// Call the main function with our test args
			os.Args = append([]string{"cmd"}, args...)
			main()

			// For now, we'll just check if files exist
			// TODO: Add proper error checking based on tt.expectError
			// TODO: Add output file validation (format, content)
			if !tt.expectError {
				if _, err := os.Stat(tt.outputPath); os.IsNotExist(err) {
					t.Errorf("Output file was not created: %s", tt.outputPath)
				}
			} else {
				if _, err := os.Stat(tt.outputPath); err == nil {
					t.Errorf("Output file was created but an error was expected: %s", tt.outputPath)
				}
			}

			// Basic input file check (should always exist for a valid test case)
			if _, err := os.Stat(tt.inputPath); os.IsNotExist(err) {
				t.Fatalf("Input file doesn't exist: %s", tt.inputPath)
			}
		})
	}
} 