package model

import "time"

type Analysis struct {
	BestNumbers       []int     `json:"best_numbers"`
	BestScores        []float64 `json:"best_scores"`
	TotalProcessed    int       `json:"total_processed"`
	GroupDaysAnalyzed int       `json:"days_analyzed"`
	ExecutionTime     string    `json:"execution_time"`
	Timestamp         time.Time `json:"timestamp"`
	UnplayedCount     int       `json:"unplayed_count"`
}

type AnalysisParams struct {
	MaxIterations int `json:"max_iterations"`
	TopNumbers    int `json:"top_numbers"`
}

type Statistics struct {
	TotalResults   int                  `json:"total_results"`
	DateRange      DateRange            `json:"date_range"`
	DigitFreq      *DigitFrequency      `json:"digit_frequency"`
	TwoDigitFreq   *TwoDigitFrequency   `json:"two_digit_frequency"`
	ThreeDigitFreq *ThreeDigitFrequency `json:"three_digit_frequency"`
	FourDigitFreq  *FourDigitFrequency  `json:"four_digit_frequency"`
}

type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type RepetitionStatistics struct {
	DigitRep      [][]float64 `json:"digit_repetition"`
	TwoDigitRep   [][]float64 `json:"two_digit_repetition"`
	ThreeDigitRep [][]float64 `json:"three_digit_repetition"`
	FourDigitRep  []float64   `json:"four_digit_repetition"`
}
