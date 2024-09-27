package excel_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/micamics/extracter/excel"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateHTTPHandler(t *testing.T) {
	var (
		svc     = new(serviceMock)
		loggr   = log.NewLogfmtLogger(os.Stderr)
		handler = excel.CreateHTTPHandler(svc, loggr)
		svr     = httptest.NewServer(handler)

		ctx      = mock.Anything
		fileMock = mock.Anything
	)

	t.Cleanup(func() {
		svc.AssertExpectations(t)
	})

	tests := []struct {
		scenario       string
		file           func(t *testing.T) (*bytes.Buffer, *multipart.Writer)
		mocks          func()
		wantStatusCode int
	}{
		{
			scenario: "valid excel file - no errors",
			file: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				return createExcelFileRequestData(t)
			},
			mocks: func() {
				svc.On("ProcessFile", ctx, fileMock).Return(nil).Once()

			},
			wantStatusCode: http.StatusOK,
		},
		{
			scenario: "valid excel file - extract data error from service",
			file: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				return createExcelFileRequestData(t)
			},
			mocks: func() {
				svc.On("ProcessFile", ctx, fileMock).Return(excel.ErrExtractingData).Once()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			scenario: "invalid file - not an excel file",
			file: func(t *testing.T) (*bytes.Buffer, *multipart.Writer) {
				return createCSVFileRequestData(t)
			},
			mocks: func() {
				svc.On("ProcessFile", ctx, fileMock).Return(excel.ErrInvalidFileType).Once()
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.scenario, func(t *testing.T) {
			tc.mocks()

			body, writer := tc.file(t)
			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/file/", svr.URL),
				body,
			)
			require.NoError(t, err)

			req.Header.Add("Content-Type", writer.FormDataContentType())

			client := http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "unable to do http request")
			defer resp.Body.Close()

			require.Equal(t, tc.wantStatusCode, resp.StatusCode)
		})
	}
}
