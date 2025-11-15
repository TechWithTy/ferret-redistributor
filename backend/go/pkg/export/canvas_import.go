package export

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ImportResult contains the result of an import operation
type ImportResult struct {
	Success bool
	Message string
	Data    interface{}
}

// ImportCanvas imports a business model canvas from a file
func ImportCanvas(filePath string, format Format) (*ImportResult, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	switch format {
	case FormatJSON:
		return importFromJSON(filePath)
	default:
		return nil, fmt.Errorf("unsupported import format: %s", format)
	}
}

// importFromJSON imports a canvas from a JSON file
func importFromJSON(filePath string) (*ImportResult, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the JSON into a map for validation
	var canvas map[string]interface{}
	if err := json.Unmarshal(data, &canvas); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// TODO: Add validation for required fields
	// This is where you would validate the structure of the imported canvas
	// and convert it to your internal model

	return &ImportResult{
		Success: true,
		Message: "Successfully imported canvas",
		Data:    canvas,
	}, nil
}

// ValidateImportFile validates if a file can be imported
func ValidateImportFile(filePath string, format Format) error {
	// Check file extension
	ext := filepath.Ext(filePath)
	if ext != ".json" && format == FormatJSON {
		return fmt.Errorf("invalid file extension for JSON import: %s", ext)
	}

	// Check file size (max 10MB)
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	if info.Size() > 10<<20 { // 10MB
		return fmt.Errorf("file is too large: %d bytes (max 10MB)", info.Size())
	}

	// Check if file is readable
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// For JSON files, check if the content is valid JSON
	if format == FormatJSON {
		decoder := json.NewDecoder(file)
		var temp interface{}
		if err := decoder.Decode(&temp); err != nil && err != io.EOF {
			return fmt.Errorf("invalid JSON content: %w", err)
		}
	}

	return nil
}
