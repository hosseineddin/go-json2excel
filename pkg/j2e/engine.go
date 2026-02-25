package j2e

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/yourusername/go-json2excel/internal/converter"
	"github.com/yourusername/go-json2excel/internal/export"
)

// Engine orchestrates the zero-copy streaming pipeline.
type Engine struct {
	inReader  io.Reader
	outWriter io.Writer
	format    export.Format
}

// NewEngine initializes a new engine instance.
func NewEngine() *Engine {
	return &Engine{
		format: export.FormatXLSX,
	}
}

func (e *Engine) SetInputReader(r io.Reader) *Engine {
	e.inReader = r
	return e
}

func (e *Engine) SetOutputWriter(w io.Writer) *Engine {
	e.outWriter = w
	return e
}

func (e *Engine) SetFormat(ext string) *Engine {
	e.format = export.Format(ext)
	return e
}

// Execute runs the highly-optimized pipeline.
func (e *Engine) Execute() error {
	if e.inReader == nil || e.outWriter == nil {
		return fmt.Errorf("input and output streams must be set")
	}

	exporter, err := export.NewExporter(e.format)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(e.inReader)
	t, err := decoder.Token()
	if err != nil || t != json.Delim('[') {
		return fmt.Errorf("expected JSON array at the root")
	}

	recordChan := make(chan map[string]interface{}, 2000)
	errChan := make(chan error, 1)

	// Goroutine 1: JSON Decoder
	go func() {
		defer close(recordChan)
		for decoder.More() {
			var record map[string]interface{}
			if err := decoder.Decode(&record); err != nil {
				errChan <- fmt.Errorf("decode error: %w", err)
				return
			}
			recordChan <- record
		}
	}()

	var headers []string
	headerMap := make(map[string]int)
	isFirstRow := true

	// Memory Pool to avoid massive Garbage Collection allocations
	rowPool := sync.Pool{
		New: func() interface{} {
			return make([]interface{}, 0, 50) // Pre-allocate capacity
		},
	}

	// Goroutine 2: Data Processing & Exporting (Main)
	for record := range recordChan {
		if isFirstRow {
			for k := range record {
				headers = append(headers, k)
				headerMap[k] = len(headers) - 1
			}
			if err := exporter.Init(headers); err != nil {
				return err
			}
			isFirstRow = false
		}

		// Fetch slice from pool to reduce heap allocations
		row := rowPool.Get().([]interface{})
		if cap(row) < len(headers) {
			row = make([]interface{}, len(headers))
		} else {
			row = row[:len(headers)]
		}

		// Map JSON keys to correct column indexes safely
		for k, v := range record {
			if idx, ok := headerMap[k]; ok {
				row[idx] = converter.ConvertDynamic(v)
			}
		}

		if err := exporter.WriteRow(row); err != nil {
			return err
		}

		// Return slice to pool
		for i := range row {
			row[i] = nil // Prevent memory leaks
		}
		rowPool.Put(row)
	}

	select {
	case err := <-errChan:
		return err
	default:
	}

	return exporter.WriteTo(e.outWriter)
}
