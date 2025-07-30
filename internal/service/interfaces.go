package service

import (
	"context"
	"time"

	"lottery-analyzer/internal/model"
)

// ScrapperService define las operaciones de scrapping de datos
type ScrapperService interface {
	ScrapingFromLastDate(ctx context.Context) error
	ScrapingDateRange(ctx context.Context, startDate, endDate time.Time) error
	LastScrapedDate(ctx context.Context) (*time.Time, error)
}

// ProcessorService define las operaciones de análisis y procesamiento
type ProcessorService interface {
	ProcessAnalysis(ctx context.Context) (*model.Analysis, error)
	BestNumbers(ctx context.Context, limit int) ([]int, []float64, error)
	UnplayedNumbers(ctx context.Context) (int, error)
	CalculateProbability(ctx context.Context, number int) (float64, error)
	Statistics(ctx context.Context) (*model.Statistics, error)
}
