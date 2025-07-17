package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"lottery-analyzer/internal/model"
	"lottery-analyzer/internal/repository"
)

type scrapperService struct {
	resultRepo repository.ResultRepository
	client     *http.Client
}

func NewScrapperService(resultRepo repository.ResultRepository) ScrapperService {
	return &scrapperService{
		resultRepo: resultRepo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *scrapperService) LastScrapedDate(ctx context.Context) (*time.Time, error) {
	result, err := s.resultRepo.LastResult(ctx)
	if err != nil || result == nil {
		return nil, err
	}

	date, err := time.Parse("02/01/2006", result.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	return &date, nil
}

func (s *scrapperService) ScrapingFromLastDate(ctx context.Context) error {
	startDate, _ := s.LastScrapedDate(ctx)

	if startDate != nil {
		tmp := startDate.AddDate(0, 0, 1) // empieza el scrapping desde el día siguiente al último hecho
		startDate = &tmp
	} else {
		tmp, _ := time.Parse("02/01/2006", "02/02/2008")
		startDate = &tmp // si no hay scrapping previo, empieza desde una fecha fija
	}

	endDate := time.Now()

	for date := *startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		dateStr := date.Format("02/01/2006")
		if err := s.scrapeDate(ctx, dateStr); err != nil {
			// Log error but continue with next date
			fmt.Printf("Failed to scrape date %s: %v\n", dateStr, err)
		}
	}

	return nil
}

func (s *scrapperService) scrapeDate(ctx context.Context, dateStr string) error {
	url := fmt.Sprintf("https://resultadodelaloteria.com/ws/services.asmx/getResultado?sFecha=%s&idLoteria=21&valueCaptcha=kZyAcju1QZE5sNoRHMohIg==&txtValueCaptcha=DMNT", dateStr)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var body strings.Builder
	if _, err = io.Copy(&body, resp.Body); err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	responseText := body.String()

	// Check if no results found
	if strings.Contains(responseText, "No se han encontrado resultados") {
		return nil // Skip, no data for this date
	}

	result, err := s.cleanData(responseText)
	if err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	result.Date = dateStr
	return s.resultRepo.Create(ctx, result)
}

func (s *scrapperService) cleanData(responseText string) (*model.Result, error) {
	index := strings.Index(responseText, "---")
	if index == -1 {
		return nil, fmt.Errorf("invalid response format")
	}

	if index < 4 {
		return nil, fmt.Errorf("insufficient data before delimiter")
	}

	// Extract number (4 digits before ---)
	numberStr := responseText[index-4 : index]

	// Extract sign
	indexEndSign := strings.Index(responseText[index+4:], "<")
	if indexEndSign == -1 {
		return nil, fmt.Errorf("sign end not found")
	}

	sign := strings.ToLower(strings.Replace(responseText[index+3:index+3+indexEndSign], "-", "", -1))

	// Convert sign to letter
	letterSign := s.convertSign(sign)

	// Parse digits
	if len(numberStr) != 4 {
		return nil, fmt.Errorf("invalid number format: %s", numberStr)
	}

	first, _ := strconv.Atoi(string(numberStr[0]))
	second, _ := strconv.Atoi(string(numberStr[1]))
	third, _ := strconv.Atoi(string(numberStr[2]))
	fourth, _ := strconv.Atoi(string(numberStr[3]))

	return &model.Result{
		First:  first,
		Second: second,
		Third:  third,
		Fourth: fourth,
		Sign:   letterSign,
	}, nil
}

func (s *scrapperService) convertSign(sign string) string {
	signMap := map[string]string{
		"acuario": "A", "acurio": "A",
		"piscis":  "B",
		"aries":   "C",
		"tauro":   "D",
		"geminis": "E", "géminis": "E",
		"cancer": "F", "cáncer": "F",
		"leo":       "G",
		"virgo":     "H",
		"libra":     "I",
		"escorpion": "J", "escorpio": "J", "escorpión": "J",
		"sagitario":   "K",
		"capricornio": "L",
	}

	if letter, ok := signMap[sign]; ok {
		return letter
	}
	return "Z"
}

func (s *scrapperService) ScrapingDateRange(ctx context.Context, startDate, endDate time.Time) error {
	for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		dateStr := date.Format("02/01/2006")
		if err := s.scrapeDate(ctx, dateStr); err != nil {
			return fmt.Errorf("failed to scrape date %s: %w", dateStr, err)
		}
	}
	return nil
}
