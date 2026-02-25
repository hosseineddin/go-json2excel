package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
)

type csvExporter struct {
	buffer *bytes.Buffer
	writer *csv.Writer
}

func newCSVExporter() *csvExporter {
	buf := new(bytes.Buffer)
	return &csvExporter{
		buffer: buf,
		writer: csv.NewWriter(buf),
	}
}

func (c *csvExporter) Init(headers []string) error {
	return c.writer.Write(headers)
}

func (c *csvExporter) WriteRow(row []interface{}) error {
	stringRow := make([]string, len(row))
	for i, val := range row {
		stringRow[i] = fmt.Sprintf("%v", val) // Fast string coercion
	}
	return c.writer.Write(stringRow)
}

func (c *csvExporter) ExportTo(w io.Writer) error {
	c.writer.Flush()
	if err := c.writer.Error(); err != nil {
		return err
	}
	_, err := w.Write(c.buffer.Bytes())
	return err
}
