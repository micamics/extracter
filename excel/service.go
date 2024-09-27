package excel

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/micamics/extracter/models"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	ProcessFile(ctx context.Context, f *models.File) error
}

func NewService() Service {
	return &service{}
}

type service struct{}

func (s *service) ProcessFile(ctx context.Context, f *models.File) error {
	excelFile, err := validateFile(f)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	for _, name := range excelFile.GetSheetMap() {

		data, err := extractDataFromSheet(excelFile, name)
		if err != nil {
			return fmt.Errorf("extract data error: %w", err)
		}

		for _, v := range data {
			if err := processData(v); err != nil {
				return fmt.Errorf("processing data error: %w", err)
			}
		}

	}

	return nil
}

func validateFile(f *models.File) (*excelize.File, error) {
	excelFile, err := excelize.OpenReader(f.Reader)
	if err != nil {
		slog.Info("open reader error: ", "error", err)
		return nil, fmt.Errorf("%w: %v", ErrInvalidFileType, ErrNotExcelType)
	}

	defer func() {
		if err := excelFile.Close(); err != nil {
			slog.Info("unable to close excel file: ", "error", err.Error())
		}
	}()

	return excelFile, nil
}

func extractDataFromSheet(f *excelize.File, sheetName string) (map[string][]string, error) {
	cols, err := f.GetCols(sheetName)
	if err != nil {
		return nil, fmt.Errorf("get column error: %w", ErrExtractingData)
	}

	dataMap := map[string][]string{}
	for _, col := range cols {
		if len(col) < 1 {
			continue
		}

		key := col[0]
		val := col[1:]

		if key == "" || len(val) == 0 {
			continue
		}

		dataMap[key] = append(dataMap[key], val...)
	}

	return dataMap, nil
}

func processData(data []string) error {
	slices.Sort(data)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("extracted_data.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Info("unable to close file", "error", err)
		}
	}()

	_, err = f.WriteString(fmt.Sprintf("%v\n", strings.Join(data, ",")))
	if err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}

	return nil
}
