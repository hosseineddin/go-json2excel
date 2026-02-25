package export

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

const maxRowsPerSheet = 1000000

type xlsxExporter struct {
	file         *excelize.File
	streamWriter *excelize.StreamWriter
	sheetIndex   int
	currentRow   int
	headers      []string
}

func newXLSXExporter() *xlsxExporter {
	return &xlsxExporter{
		file:       excelize.NewFile(),
		sheetIndex: 1,
		currentRow: 1,
	}
}

func (x *xlsxExporter) Init(headers []string) error {
	x.headers = headers
	sheetName := fmt.Sprintf("Sheet%d", x.sheetIndex)

	if x.sheetIndex > 1 {
		x.file.NewSheet(sheetName)
	} else {
		x.file.SetSheetName("Sheet1", sheetName)
	}

	sw, err := x.file.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}
	x.streamWriter = sw

	headerRow := make([]interface{}, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}

	cell, _ := excelize.CoordinatesToCellName(1, 1)
	if err := x.streamWriter.SetRow(cell, headerRow); err != nil {
		return err
	}
	x.currentRow = 2
	return nil
}

func (x *xlsxExporter) WriteRow(row []interface{}) error {
	if x.currentRow > maxRowsPerSheet {
		if err := x.streamWriter.Flush(); err != nil {
			return err
		}
		x.sheetIndex++
		x.currentRow = 1
		if err := x.Init(x.headers); err != nil {
			return err
		}
	}

	cell, _ := excelize.CoordinatesToCellName(1, x.currentRow)
	if err := x.streamWriter.SetRow(cell, row); err != nil {
		return err
	}
	x.currentRow++
	return nil
}

func (x *xlsxExporter) WriteTo(w io.Writer) error {
	if err := x.streamWriter.Flush(); err != nil {
		return err
	}
	return x.file.Write(w)
}
