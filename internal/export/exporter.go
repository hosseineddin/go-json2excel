package export

import (
	"fmt"
	"io"
)

// Exporter defines the contract for streaming outputs.
type Exporter interface {
	Init(headers []string) error
	WriteRow(row []interface{}) error
	ExportTo(w io.Writer) error
}

// Format defines the requested output extension.
type Format string

const (
	FormatXLSX Format = "xlsx"
	FormatCSV  Format = "csv"
)

// NewExporter is a Factory returning the correct optimized writer.
func NewExporter(format Format) (Exporter, error) {
	switch format {
	case FormatXLSX:
		return newXLSXExporter(), nil
	case FormatCSV:
		return newCSVExporter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
