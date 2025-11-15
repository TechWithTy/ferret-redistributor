package export

// Format represents a supported file format for both import and export operations
type Format string

const (
    // FormatJSON represents JSON file format
    FormatJSON Format = "json"
    // FormatPDF represents PDF file format (export only)
    FormatPDF Format = "pdf"
    // FormatPNG represents PNG file format (export only)
    FormatPNG Format = "png"
)

// SupportedImportFormats returns a list of supported import formats
func SupportedImportFormats() []Format {
    return []Format{FormatJSON}
}

// SupportedExportFormats returns a list of supported export formats
func SupportedExportFormats() []Format {
    return []Format{FormatJSON, FormatPDF, FormatPNG}
}
