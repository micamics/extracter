package excel

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrUploading       = errors.New("unable to upload")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrNotExcelType    = errors.New("not an excel file")
	ErrExtractingData  = errors.New("unable to extract data")

	internalServerError = "There are few issues that we need to fix. " +
		"We are actively working on it. Thank you for your patience."
)

func prepareHTTPError(err error) (code int, errMsg string) {
	badErrs := []error{
		ErrUploading,
		ErrInvalidFileType,
		ErrNotExcelType,
	}

	code, errMsg = prepareBadRequestError(err, badErrs)
	if code != 0 && errMsg != "" {
		return code, errMsg
	}

	return http.StatusInternalServerError, internalServerError
}

func prepareBadRequestError(err error, badRequestErrs []error) (code int, errMsg string) {
	for _, badErr := range badRequestErrs {
		if strings.Contains(err.Error(), badErr.Error()) {
			code = http.StatusBadRequest
			errMsg = badErr.Error()

			break
		}
	}

	return code, errMsg
}
