package excel_test

import (
	"context"

	"github.com/micamics/extracter/models"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) ProcessFile(
	ctx context.Context,
	f *models.File,
) error {
	return s.Called(ctx, f).Error(0)
}
