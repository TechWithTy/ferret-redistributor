package export

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

// ExportCanvas exports a business model canvas to the specified format

// ExportCanvas exports a business model canvas to the specified format
func ExportCanvas(canvas interface{}, format Format, outputPath string) error {
	switch format {
	case FormatPDF:
		return exportToPDF(canvas, outputPath)
	case FormatPNG:
		return exportToPNG(canvas, outputPath)
	case FormatJSON:
		return exportToJSON(canvas, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToPDF exports the canvas to a PDF file
func exportToPDF(canvas interface{}, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Set font and colors
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	
	// Add title
	pdf.Cell(40, 10, "Business Model Canvas")
	pdf.Ln(20)
	
	// TODO: Add canvas content rendering here
	// This is a simplified example - you'll need to implement the actual rendering
	// of your canvas components
	
	// Create output directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Save the PDF
	return pdf.OutputFileAndClose(outputPath)
}

// exportToPNG exports the canvas to a PNG file
func exportToPNG(canvas interface{}, outputPath string) error {
	// Create a simple image (replace with actual canvas rendering)
	img := image.NewRGBA(image.Rect(0, 0, 1200, 800))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	
	// TODO: Add canvas content rendering here
	// This is a simplified example - you'll need to implement the actual rendering
	// of your canvas components using the draw package or a more advanced library
	
	// Create output directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Save the image
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()
	
	return png.Encode(f, img)
}

// exportToJSON exports the canvas to a JSON file
func exportToJSON(canvas interface{}, outputPath string) error {
	// Create output directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(canvas, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal canvas to JSON: %w", err)
	}
	
	// Write to file
	return os.WriteFile(outputPath, data, 0644)
}
