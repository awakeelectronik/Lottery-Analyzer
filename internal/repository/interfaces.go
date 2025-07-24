package repository

import (
	"context"
	"time"

	"lottery-analyzer/internal/model"
)

// ResultRepository define las operaciones de acceso a datos para Result
type ResultRepository interface {
	Create(ctx context.Context, result *model.Result) error
	LastResult(ctx context.Context) (*model.Result, error)
	OneDigit(ctx context.Context, cal time.Time, position string) ([]*model.DigitCount, error)
	TwoDigit(ctx context.Context, cal time.Time, position1 string, position2 string) ([]*model.TwoDigitCount, error)
	CreateBatch(ctx context.Context, results []*model.Result) error
	ID(ctx context.Context, id int) (*model.Result, error)
	Date(ctx context.Context, date string) ([]*model.Result, error)
	LastNResults(ctx context.Context, limit int) ([]*model.Result, error)
	BetweenDates(ctx context.Context, startDate, endDate time.Time) ([]*model.Result, error)
	AfterDate(ctx context.Context, date time.Time) ([]*model.Result, error)
	Update(ctx context.Context, result *model.Result) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, date string) (bool, error)
	Count(ctx context.Context) (int, error)
	CountBetweenDates(ctx context.Context, startDate, endDate time.Time) (int, error)
	AllPlayedNumbers(ctx context.Context) ([]string, error)
}
