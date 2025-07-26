package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"lottery-analyzer/internal/model"
	"lottery-analyzer/pkg/utils"
)

type resultRepository struct {
	db *sql.DB
}

func NewResultRepository(db *sql.DB) ResultRepository {
	return &resultRepository{db: db}
}

func (r *resultRepository) Create(ctx context.Context, result *model.Result) error {
	query := `INSERT INTO result (version, date, first, second, third, fourth, sign) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		result.Version, result.Date, result.First, result.Second,
		result.Third, result.Fourth, result.Sign)

	return err
}

func (r *resultRepository) LastResult(ctx context.Context) (*model.Result, error) {
	query := `SELECT id, version, date, first, second, third, fourth, sign 
              FROM result ORDER BY id DESC LIMIT 1`

	var result model.Result
	err := r.db.QueryRowContext(ctx, query).Scan(
		&result.ID, &result.Version, &result.Date,
		&result.First, &result.Second, &result.Third, &result.Fourth, &result.Sign)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (r *resultRepository) OneDigit(ctx context.Context, cal time.Time, position string) ([]*model.DigitCount, error) {

	escapedCol := utils.EscapeIdentifier(position) // Escapa para evitar SQL injection

	query := fmt.Sprintf(`SELECT COUNT(%s) AS repetition, %s 
                          FROM result 
                          WHERE STR_TO_DATE(date, ?) > ? 
                          GROUP BY %s`,
		escapedCol, escapedCol, escapedCol)

	// Ejecuta la query con placeholders solo para valores
	rows, err := r.db.QueryContext(ctx, query, "%d/%m/%Y", cal)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err) // Manejo de errores con wrapping
	}
	defer rows.Close()

	var counts []*model.DigitCount
	for rows.Next() {
		var dc model.DigitCount
		if err := rows.Scan(&dc.Count, &dc.Digit); err != nil {
			return nil, err
		}
		counts = append(counts, &dc)
	}
	return counts, rows.Err()
}

func (r *resultRepository) TwoDigit(ctx context.Context, cal time.Time, position1 string, position2 string) ([]*model.TwoDigitCount, error) {

	escapedCol1 := utils.EscapeIdentifier(position1) // Escapa para evitar SQL injection
	escapedCol2 := utils.EscapeIdentifier(position2) // Escapa para evitar SQL injection

	query := fmt.Sprintf(`SELECT COUNT(%s) AS repetition, %s, %s 
                          FROM result 
                          WHERE STR_TO_DATE(date, ?) > ? 
                          GROUP BY %s, %s`,
		escapedCol1, escapedCol1, escapedCol2, escapedCol1, escapedCol2)

	// Ejecuta la query con placeholders solo para valores
	rows, err := r.db.QueryContext(ctx, query, "%d/%m/%Y", cal)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err) // Manejo de errores con wrapping
	}
	defer rows.Close()

	var counts []*model.TwoDigitCount
	for rows.Next() {
		var dc model.TwoDigitCount
		if err := rows.Scan(&dc.Count, &dc.FirstDigit, &dc.SecondDigit); err != nil {
			return nil, err
		}
		counts = append(counts, &dc)
	}
	return counts, rows.Err()
}

func (r *resultRepository) ThreeDigit(ctx context.Context, cal time.Time, position1 string, position2 string, position3 string) ([]*model.ThreeDigitCount, error) {

	escapedCol1 := utils.EscapeIdentifier(position1) // Escapa para evitar SQL injection
	escapedCol2 := utils.EscapeIdentifier(position2) // Escapa para evitar SQL injection
	escapedCol3 := utils.EscapeIdentifier(position3) // Escapa para evitar SQL injection

	query := fmt.Sprintf(`SELECT COUNT(%s) AS repetition, %s, %s, %s
                          FROM result 
                          WHERE STR_TO_DATE(date, ?) > ? 
                          GROUP BY %s, %s, %s`,
		escapedCol1, escapedCol1, escapedCol2, escapedCol3, escapedCol1, escapedCol2, escapedCol3)

	// Ejecuta la query con placeholders solo para valores
	rows, err := r.db.QueryContext(ctx, query, "%d/%m/%Y", cal)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err) // Manejo de errores con wrapping
	}
	defer rows.Close()

	var counts []*model.ThreeDigitCount
	for rows.Next() {
		var dc model.ThreeDigitCount
		if err := rows.Scan(&dc.Count, &dc.FirstDigit, &dc.SecondDigit, &dc.ThirdDigit); err != nil {
			return nil, err
		}
		counts = append(counts, &dc)
	}
	return counts, rows.Err()
}

func (r *resultRepository) FourthDigit(ctx context.Context, cal time.Time) ([]*model.FourDigitCount, error) {
	query := `SELECT COUNT(first) AS repetition, first, second, third, fourth 
                          FROM result 
                          WHERE STR_TO_DATE(date, ?) > ? 
                          GROUP BY first, second, third, fourth`

	rows, err := r.db.QueryContext(ctx, query, "%d/%m/%Y", cal)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var counts []*model.FourDigitCount
	for rows.Next() {
		var first, second, third, fourth int
		dc := model.FourDigitCount{}
		if err := rows.Scan(&dc.Count, &first, &second, &third, &fourth); err != nil {
			return nil, err
		}
		dc.Number = first*1000 + second*100 + third*10 + fourth
		counts = append(counts, &dc)
	}
	return counts, rows.Err()
}

func (r *resultRepository) CreateBatch(ctx context.Context, results []*model.Result) error {
	if len(results) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO result (version, date, first, second, third, fourth, sign) 
         VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, result := range results {
		_, err := stmt.ExecContext(ctx,
			result.Version, result.Date, result.First, result.Second,
			result.Third, result.Fourth, result.Sign)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *resultRepository) ID(ctx context.Context, id int) (*model.Result, error) {
	query := `SELECT id, version, date, first, second, third, fourth, sign 
              FROM result WHERE id = ?`

	var result model.Result
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.ID, &result.Version, &result.Date,
		&result.First, &result.Second, &result.Third, &result.Fourth, &result.Sign)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (r *resultRepository) Date(ctx context.Context, date string) ([]*model.Result, error) {
	query := `SELECT id, version, date, first, second, third, fourth, sign 
              FROM result WHERE date = ?`

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Result
	for rows.Next() {
		var result model.Result
		if err := rows.Scan(&result.ID, &result.Version, &result.Date,
			&result.First, &result.Second, &result.Third, &result.Fourth, &result.Sign); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, rows.Err()
}

func (r *resultRepository) LastNResults(ctx context.Context, limit int) ([]*model.Result, error) {
	query := `SELECT id, version, date, first, second, third, fourth, sign 
              FROM result ORDER BY id DESC LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Result
	for rows.Next() {
		var result model.Result
		if err := rows.Scan(&result.ID, &result.Version, &result.Date,
			&result.First, &result.Second, &result.Third, &result.Fourth, &result.Sign); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, rows.Err()
}

func (r *resultRepository) BetweenDates(ctx context.Context, startDate, endDate time.Time) ([]*model.Result, error) {
	query := `SELECT id, version, date, first, second, third, fourth, sign 
              FROM result WHERE STR_TO_DATE(date, '%d/%m/%Y') BETWEEN ? AND ?
              ORDER BY STR_TO_DATE(date, '%d/%m/%Y')`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Result
	for rows.Next() {
		var result model.Result
		if err := rows.Scan(&result.ID, &result.Version, &result.Date,
			&result.First, &result.Second, &result.Third, &result.Fourth, &result.Sign); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, rows.Err()
}

func (r *resultRepository) AfterDate(ctx context.Context, date time.Time) ([]*model.Result, error) {
	query := `SELECT id, version, date, first, second, third, fourth, sign 
              FROM result WHERE STR_TO_DATE(date, '%d/%m/%Y') > ?
              ORDER BY STR_TO_DATE(date, '%d/%m/%Y')`

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Result
	for rows.Next() {
		var result model.Result
		if err := rows.Scan(&result.ID, &result.Version, &result.Date,
			&result.First, &result.Second, &result.Third, &result.Fourth, &result.Sign); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, rows.Err()
}

func (r *resultRepository) Update(ctx context.Context, result *model.Result) error {
	query := `UPDATE result SET version = ?, date = ?, first = ?, second = ?, 
              third = ?, fourth = ?, sign = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		result.Version, result.Date, result.First, result.Second,
		result.Third, result.Fourth, result.Sign, result.ID)

	return err
}

func (r *resultRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM result WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *resultRepository) Exists(ctx context.Context, date string) (bool, error) {
	query := `SELECT COUNT(*) FROM result WHERE date = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, date).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *resultRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM result`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *resultRepository) CountBetweenDates(ctx context.Context, startDate, endDate time.Time) (int, error) {
	query := `SELECT COUNT(*) FROM result 
              WHERE STR_TO_DATE(date, '%d/%m/%Y') BETWEEN ? AND ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, startDate, endDate).Scan(&count)
	return count, err
}

func (r *resultRepository) AllPlayedNumbers(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT CONCAT(LPAD(first, 1, '0'), LPAD(second, 1, '0'), 
              LPAD(third, 1, '0'), LPAD(fourth, 1, '0')) as number FROM result`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var numbers []string
	for rows.Next() {
		var number string
		if err := rows.Scan(&number); err != nil {
			return nil, err
		}
		numbers = append(numbers, number)
	}

	return numbers, rows.Err()
}
