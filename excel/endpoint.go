package excel

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/micamics/extracter/models"
)

// CreateServerEndpoints initializes new endpoints for the upload picture service.
func CreateServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ProcessFileEndpoint: MakeProcessFileEndpoint(s),
	}
}

type Endpoints struct {
	ProcessFileEndpoint endpoint.Endpoint
}

type processFileRequest struct {
	File *models.File
}

type processFileResponse struct {
	Err error `json:"error,omitempty"`
}

func MakeProcessFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(processFileRequest)
		err := s.ProcessFile(ctx, req.File)

		return processFileResponse{
			Err: err,
		}, nil
	}
}
