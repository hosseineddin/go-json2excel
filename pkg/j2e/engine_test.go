package j2e

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

// generateMockJSON streams mock JSON records directly to an io.Writer
func generateMockJSON(w io.Writer, records int) {
	w.Write([]byte("[\n"))
	for i := 0; i < records; i++ {
		row := fmt.Sprintf(`{"id":%d,"name":"User_%d","score":%f,"active":true}`, i, i, float64(i)*1.5)
		if i < records-1 {
			row += ",\n"
		} else {
			row += "\n"
		}
		w.Write([]byte(row))
	}
	w.Write([]byte("]\n"))
}

func TestEngine_Execute_CSV(t *testing.T) {
	inBuf := new(bytes.Buffer)
	generateMockJSON(inBuf, 1000000) // 1000 records for functional test

	outBuf := new(bytes.Buffer)

	err := NewEngine().
		SetInputReader(inBuf).
		SetOutputWriter(outBuf).
		SetFormat("csv").
		Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if outBuf.Len() == 0 {
		t.Fatal("Expected output buffer to contain data, got empty")
	}
}

// Benchmark the engine with 100,000 records
// Run with: go test -bench=. -benchmem
func BenchmarkEngine_Execute_CSV_100k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer() // Pause timer while generating data
		inBuf := new(bytes.Buffer)
		generateMockJSON(inBuf, 100000)
		outBuf := new(bytes.Buffer) // Use discard in real-world pure CPU bench: io.Discard
		b.StartTimer()              // Resume timer

		err := NewEngine().
			SetInputReader(inBuf).
			SetOutputWriter(outBuf).
			SetFormat("csv").
			Execute()

		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
