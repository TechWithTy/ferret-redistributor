# Business Model Canvas - Import/Export Guide

This document outlines the import and export functionality for the Business Model Canvas.

## Table of Contents
- [Supported Formats](#supported-formats)
- [Exporting a Canvas](#exporting-a-canvas)
- [Importing a Canvas](#importing-a-canvas)
- [File Structure](#file-structure)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Supported Formats

| Format | Import | Export | Description                     |
|--------|--------|--------|---------------------------------|
| JSON   | ✅      | ✅      | Full canvas data with metadata  |
| PDF    | ❌      | ✅      | Printable document              |
| PNG    | ❌      | ✅      | Image export                    |

## Exporting a Canvas

### Code Example

```go
import "github.com/your-org/ferret-redistributor/go/pkg/export"

// Export to PDF
err := export.ExportCanvas(canvas, export.FormatPDF, "output/business_model.pdf")
if err != nil {
    log.Fatalf("Export failed: %v", err)
}

// Export to JSON
err = export.ExportCanvas(canvas, export.FormatJSON, "output/business_model.json")
```

## Importing a Canvas

### Code Example

```go
import "github.com/your-org/ferret-redistributor/go/pkg/export"

// Import from JSON
result, err := export.ImportCanvas("import/canvas.json", export.FormatJSON)
if err != nil {
    log.Fatalf("Import failed: %v", err)
}

if result.Success {
    fmt.Println("Successfully imported canvas:", result.Message)
    canvas := result.Data.(*models.BusinessModelCanvas)
    // Use the imported canvas
}
```

## File Structure

### JSON Format
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "My Business Model",
  "description": "Example business model canvas",
  "customerSegments": [
    {
      "id": "segment-1",
      "name": "Small Businesses",
      "description": "Local small business owners"
    }
  ],
  "valuePropositions": [
    {
      "id": "vp-1",
      "name": "Affordable Solution",
      "description": "Cost-effective alternative to competitors"
    }
  ]
  // Other canvas components...
}
```

## Error Handling

### Common Errors

| Error Code | Description | Resolution |
|------------|-------------|------------|
| `ERR_FILE_NOT_FOUND` | Specified file doesn't exist | Verify file path |
| `ERR_INVALID_FORMAT` | Unsupported file format | Use .json, .pdf, or .png |
| `ERR_VALIDATION` | Invalid canvas data | Check data structure |

### Handling Errors

```go
result, err := export.ImportCanvas("invalid.json", export.FormatJSON)
if err != nil {
    switch {
    case errors.Is(err, export.ErrFileNotFound):
        // Handle missing file
    case errors.Is(err, export.ErrInvalidFormat):
        // Handle invalid format
    default:
        // Handle other errors
    }
}
```

## Examples

### Export to PDF
```go
err := export.ExportCanvas(canvas, export.FormatPDF, "export.pdf")
```

### Export to PNG
```go
err := export.ExportCanvas(canvas, export.FormatPNG, "export.png")
```

### Import from JSON
```go
result, err := export.ImportCanvas("backup.json", export.FormatJSON)
if err == nil && result.Success {
    // Handle successful import
}
```

## Best Practices

1. Always validate imports before processing
2. Use absolute file paths for exports
3. Handle large exports asynchronously
4. Implement proper error handling and user feedback
5. Consider file size limits for imports

## Troubleshooting

**Issue**: Import fails with validation errors  
**Solution**: Check that all required fields are present and properly formatted

**Issue**: Export file is empty  
**Solution**: Verify write permissions and disk space

**Issue**: Imported data doesn't match expected format  
**Solution**: Validate against the schema before import
