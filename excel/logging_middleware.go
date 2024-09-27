package excel

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/micamics/extracter/models"
)

// Middleware describes a service middleware.
type Middleware func(Service) Service

func Logging(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &logging{
			next:   next,
			logger: logger,
		}
	}
}

type logging struct {
	next   Service
	logger log.Logger
}

func (mw logging) ProcessFile(
	ctx context.Context,
	f *models.File,
) (err error) {
	defer func(begin time.Time) {
		//nolint:errcheck // We don't care about the error.
		mw.logger.Log(
			"method", "ProcessFile",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return mw.next.ProcessFile(ctx, f)
}
