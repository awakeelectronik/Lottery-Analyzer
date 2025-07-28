package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"lottery-analyzer/internal/model"
	"lottery-analyzer/internal/repository"
)

type processorService struct {
	scrapperService ScrapperService
	resultRepo      repository.ResultRepository
}

func NewProcessorService(scrapper ScrapperService, resultRepo repository.ResultRepository) ProcessorService {
	return &processorService{
		scrapperService: scrapper,
		resultRepo:      resultRepo,
	}
}

func (p *processorService) ProcessAnalysis(ctx context.Context) (*model.Analysis, error) {
	start := time.Now()

	// 1. Ejecutar scrapping
	if err := p.scrapperService.ScrapingFromLastDate(ctx); err != nil {
		return nil, fmt.Errorf("scrapping failed: %w", err)
	}

	// 2. Inicializar matrices de repetición
	// eliminar después
	repetitionDigit := make([][]float64, 4)
	for i := range repetitionDigit {
		repetitionDigit[i] = make([]float64, 10)
	}

	repetitionTwoDigit := make([][]float64, 6)
	for i := range repetitionTwoDigit {
		repetitionTwoDigit[i] = make([]float64, 100)
	}

	repetitionThreeDigit := make([][]float64, 4)
	for i := range repetitionThreeDigit {
		repetitionThreeDigit[i] = make([]float64, 1000)
	}

	repetitionFourDigit := make([]float64, 10000)

	// 3. Secuencia Fibonacci para fechas
	frequencyData, fibCalcuCount, err := CalculateFrequencies(ctx, p.resultRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate frequencies: %w", err)
	}

	if frequencyData == nil {
		return nil, fmt.Errorf("frequency data is nil")
	}

	// 4. Calcular probabilidades y encontrar mejores números
	bestNumbers := make([]int, 100)
	bestScores := make([]float64, 100)

	// Inicializar con valores altos
	for i := range bestScores {
		bestScores[i] = 10000.0
		bestNumbers[i] = 1000 - i
	}

	for number := 0; number < 10000; number++ {
		score := p.calculateProbabilityResult(number, repetitionDigit, repetitionTwoDigit, repetitionThreeDigit, repetitionFourDigit)

		// Si el score es mejor que el peor de los mejores
		if bestScores[99] > score {
			bestScores[99] = score
			bestNumbers[99] = number

			// Ordenar y reorganizar
			p.sortBestNumbers(bestNumbers, bestScores)
		}
	}

	// 5. Calcular números no jugados
	unplayedCount, err := p.UnplayedNumbers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unplayed numbers: %w", err)
	}

	return &model.Analysis{
		BestNumbers:    bestNumbers,
		BestScores:     bestScores,
		TotalProcessed: 10000,
		DaysAnalyzed:   fibCalcuCount,
		ExecutionTime:  time.Since(start).String(),
		Timestamp:      time.Now(),
		UnplayedCount:  unplayedCount,
	}, nil
}

func (p *processorService) ProcessAnalysisWithParams(ctx context.Context, params *model.AnalysisParams) (*model.Analysis, error) {
	// Implementación simplificada usando ProcessAnalysis base
	return p.ProcessAnalysis(ctx)
}

func (p *processorService) BestNumbers(ctx context.Context, limit int) ([]int, []float64, error) {
	analysis, err := p.ProcessAnalysis(ctx)
	if err != nil {
		return nil, nil, err
	}

	if limit > len(analysis.BestNumbers) {
		limit = len(analysis.BestNumbers)
	}

	return analysis.BestNumbers[:limit], analysis.BestScores[:limit], nil
}

func (p *processorService) UnplayedNumbers(ctx context.Context) (int, error) {
	// Universo de números posibles (0000-9999)
	universe := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		universe[fmt.Sprintf("%04d", i)] = true
	}

	// Obtener números que ya han salido
	playedNumbers, err := p.resultRepo.AllPlayedNumbers(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get played numbers: %w", err)
	}

	// Eliminar números jugados del universo
	for _, played := range playedNumbers {
		delete(universe, played)
	}

	return len(universe), nil
}

func (p *processorService) CalculateProbability(ctx context.Context, number int) (float64, error) {
	// Implementación simplificada
	return 0.0, nil
}

func (p *processorService) Statistics(ctx context.Context) (*model.Statistics, error) {
	// Implementación básica de estadísticas
	count, err := p.resultRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &model.Statistics{
		TotalResults: count,
		DateRange: model.DateRange{
			StartDate: time.Now().AddDate(-1, 0, 0),
			EndDate:   time.Now(),
		},
	}, nil
}

// Funciones auxiliares

func (p *processorService) calculateProbabilityResult(number int, digitRep [][]float64, twoDigitRep [][]float64, threeDigitRep [][]float64, fourDigitRep []float64) float64 {
	numberStr := fmt.Sprintf("%04d", number)
	var prob float64

	// Convertir a enteros para indexación
	digits := make([]int, 4)
	for i, char := range numberStr {
		digits[i], _ = strconv.Atoi(string(char))
	}

	// Probabilidad por dígito individual
	prob += digitRep[0][digits[0]]
	prob += digitRep[1][digits[1]]
	prob += digitRep[2][digits[2]]
	prob += digitRep[3][digits[3]]

	// Probabilidad por dos dígitos
	prob += twoDigitRep[0][digits[0]*10+digits[1]] // first+second
	prob += twoDigitRep[1][digits[1]*10+digits[2]] // second+third
	prob += twoDigitRep[2][digits[2]*10+digits[3]] // third+fourth
	prob += twoDigitRep[3][digits[0]*10+digits[2]] // first+third
	prob += twoDigitRep[4][digits[0]*10+digits[3]] // first+fourth
	prob += twoDigitRep[5][digits[1]*10+digits[3]] // second+fourth

	// Probabilidad por tres dígitos
	prob += threeDigitRep[0][digits[0]*100+digits[1]*10+digits[2]] // first+second+third
	prob += threeDigitRep[1][digits[1]*100+digits[2]*10+digits[3]] // second+third+fourth
	prob += threeDigitRep[2][digits[2]*100+digits[3]*10+digits[0]] // third+fourth+first
	prob += threeDigitRep[3][digits[1]*100+digits[3]*10+digits[0]] // second+fourth+first

	return prob
}

func (p *processorService) sortBestNumbers(numbers []int, scores []float64) {
	// Crear slice de índices para ordenamiento
	indices := make([]int, len(scores))
	for i := range indices {
		indices[i] = i
	}

	// Ordenar por scores
	sort.Slice(indices, func(i, j int) bool {
		return scores[indices[i]] < scores[indices[j]]
	})

	// Reordenar arrays
	tempNumbers := make([]int, len(numbers))
	tempScores := make([]float64, len(scores))

	for i, idx := range indices {
		tempNumbers[i] = numbers[idx]
		tempScores[i] = scores[idx]
	}

	copy(numbers, tempNumbers)
	copy(scores, tempScores)
}
