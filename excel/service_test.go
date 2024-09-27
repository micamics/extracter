package excel_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/micamics/extracter/excel"
	"github.com/micamics/extracter/models"
	"github.com/stretchr/testify/require"
)

func TestService_ProcessFile(t *testing.T) {
	var (
		svc = excel.NewService()

		xcel = createTestExcelFile(t, 5)
		csv  = createTestCSVFile(t)
	)

	tests := []struct {
		scenario   string
		fileReader func() io.Reader
		err        error
	}{
		{
			scenario: "valid excel file",
			fileReader: func() io.Reader {
				buf, err := xcel.WriteToBuffer()
				require.NoError(t, err)

				return buf
			},
			err: nil,
		},
		{
			scenario: "invalid non-excel - csv file",
			fileReader: func() io.Reader {
				var buf []byte

				_, err := csv.Read(buf)
				require.NoError(t, err)

				return bytes.NewBuffer(buf)
			},
			err: excel.ErrInvalidFileType,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			requestFile := &models.File{Reader: tc.fileReader()}
			err := svc.ProcessFile(context.Background(), requestFile)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}
