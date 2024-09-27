package excel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/micamics/extracter/models"

	httptransport "github.com/go-kit/kit/transport/http"
)

// CreateHTTPHandler creates handlers for excel file processing endpoint.
func CreateHTTPHandler(svc Service, loggr log.Logger) http.Handler {
	r := mux.NewRouter()
	e := CreateServerEndpoints(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(loggr),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	r.Methods(http.MethodPost).
		Path("/file/").Handler(
		httptransport.NewServer(
			e.ProcessFileEndpoint,
			decodeProcessFileRequest,
			encodeProcessFileResponse,
			options...,
		),
	)

	maxAllowed := 10 * 1024 * 1024 // 10mb
	handler := http.MaxBytesHandler(r, int64(maxAllowed))

	return handler
}

const maxMemory = 5 * 1024 * 1024 // 5 megabytes.
func decodeProcessFileRequest(_ context.Context, r *http.Request) (request any, err error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, fmt.Errorf("unable to parse form: %v", err)
	}

	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			panic(err)
		}
	}()

	f, _, err := r.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, fmt.Errorf("%w: empty file not allowed", ErrUploading)
		}
		return nil, fmt.Errorf("%w: %v", ErrUploading, err)
	}
	defer f.Close()

	return processFileRequest{
		File: &models.File{
			Reader: f,
		},
	}, nil
}

func encodeProcessFileResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res, ok := response.(processFileResponse)
	if ok && res.Err != nil {
		return res.Err
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	statusCode, msg := prepareHTTPError(err)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
	}); err != nil {
		slog.Error("unable to encode error response: ", "error", err.Error())
	}
}
