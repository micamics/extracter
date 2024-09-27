package excel_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func createTestExcelFile(t *testing.T, numColumns int) *excelize.File {
	t.Helper()

	file := excelize.NewFile()

	// We'll create a test file with two sheets.
	sheetNum := 2
	numRows := 5
	startCol := 2

	numColumns += startCol

	for s := 1; s <= sheetNum; s++ {
		sheetName := fmt.Sprintf("Sheet%d", s)

		//nolint:errcheck // We don't care about the error.
		file.NewSheet(sheetName)

		for i := 1; i <= numRows; i++ {
			// We're doing this so the first column is empty.
			// This is to test if empty column is not evaluated.
			for j := startCol; j <= numColumns; j++ {
				value := fmt.Sprintf("%v-data-%v%v", sheetName, j, i)
				t.Logf("val of i: %v", i)
				if i == 1 {
					value = fmt.Sprintf("Header %v", j-i)
				}

				colName, err := excelize.ColumnNumberToName(j)
				if err != nil {
					require.NoError(t, err)
				}

				cellAxis := fmt.Sprintf("%v%v", colName, i)
				err = file.SetCellValue(sheetName, cellAxis, value)
				require.NoError(t, err)
			}
		}
	}

	//nolint:errcheck // We don't care about the error.
	file.SaveAs("test_excel.xlsx")

	t.Cleanup(func() {
		file.Close()
	})

	return file
}

func createExcelFileRequestData(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file := createTestExcelFile(t, 5)
	buf, err := file.WriteToBuffer()
	require.NoError(t, err)

	part, err := writer.CreateFormFile("file", "test_excel.xlsx")
	require.NoError(t, err, "unable to create form file")

	_, err = io.Copy(part, buf)
	require.NoError(t, err, "unable to copy file")

	writer.Close()

	return payload, writer
}

func createTestCSVFile(t *testing.T) *os.File {
	t.Helper()

	file, err := os.Create("file.csv")
	require.NoError(t, err)

	_, err = file.WriteString("The quick brown fox jumps over the lazy dog")
	require.NoError(t, err)

	t.Cleanup(func() {
		file.Close()
	})

	return file
}

func createCSVFileRequestData(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	csvFile := createTestCSVFile(t)

	file, err := os.Open(csvFile.Name())
	require.NoError(t, err, "unable to open file")

	part, err := writer.CreateFormFile("file", csvFile.Name())
	require.NoError(t, err, "unable to create form file")

	_, err = io.Copy(part, file)
	require.NoError(t, err, "unable to copy file")

	writer.Close()

	return payload, writer
}
