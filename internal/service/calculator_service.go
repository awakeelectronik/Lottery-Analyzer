package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"lottery-analyzer/internal/model"
	"lottery-analyzer/internal/repository"
)

// CalculateFrequencies calculates the frequencies of digits, two-digit, three-digit, and four-digit numbers
// Ya después se pueden calcular las probabilidades de un número de cuatro cifras usando las probabilidades de sus combinaciones.

func CalculateFrequencies(ctx context.Context, resultRepo repository.ResultRepository) (int, error) {

	before, actual, frequenciesProcessed := 1, 1, 0

	// contenedor de las probabilidades por dígito
	digit := &model.DigitFrequency{
		Position1: make([]float64, 10),
		Position2: make([]float64, 10),
		Position3: make([]float64, 10),
		Position4: make([]float64, 10),
	}

	twoDigit := &model.TwoDigitFrequency{
		FirstSecond:  make([]float64, 100),
		FirstThird:   make([]float64, 100),
		FirstFourth:  make([]float64, 100),
		SecondThird:  make([]float64, 100),
		SecondFourth: make([]float64, 100),
		ThirdFourth:  make([]float64, 100),
	}

	threeDigit := &model.ThreeDigitFrequency{
		FirstSecondThird:  make([]float64, 1000),
		FirstSecondFourth: make([]float64, 1000),
		FirstThirdFourth:  make([]float64, 1000),
		SecondThirdFourth: make([]float64, 1000),
	}

	fourDigit := &model.FourDigitFrequency{
		Complete: make([]float64, 10000),
	}

	// to-do este 5000 hay que volverlo automático y pensarse la lógica que está usando para optimizarla.
	for actual < 5000 {

		date := time.Now().AddDate(0, 0, -(actual + 7))

		if err := digitFrequencies(ctx, date, resultRepo, digit); err != nil {
			return 0, fmt.Errorf("digit frequency query failed: %w", err)
		}

		if err := twoDigitFrequencies(ctx, date, resultRepo, twoDigit); err != nil {
			return 0, fmt.Errorf("two digit frequency query failed: %w", err)
		}

		if err := threeDigitFrequencies(ctx, date, resultRepo, threeDigit); err != nil {
			return 0, fmt.Errorf("three digit frequency query failed: %w", err)
		}

		if err := fourDigitFrequencies(ctx, date, resultRepo, fourDigit); err != nil {
			return 0, fmt.Errorf("four digit frequency query failed: %w", err)
		}

		before, actual = actual, before+actual
		//fmt.Printf("actual: %d \n", actual)
		frequenciesProcessed++
	}

	for i, v := range fourDigit.Complete {
		if v > 100 {
			fmt.Printf("Digit %d frequency is %f \n", i, v)

		}
	}

	return frequenciesProcessed, nil
}

// Implementar métodos de consulta de frecuencias por commbinación de uno, dos y tres dígitos, también con todos los dígitos.

func digitFrequencies(ctx context.Context, fromDate time.Time, resultRepo repository.ResultRepository, digitCount *model.DigitFrequency) error {
	factors := []float64{1.0, 1.9, 2.69, 2.69} //susceptible de ser parámetros si hago computación evolutiva, por eso lo dejo como variable to-do-ia

	counts, err := resultRepo.OneDigit(ctx, fromDate, "first")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by digit-first: %w", err)
	} else {
		digitCount.Position1 = sumProbDigit(counts, digitCount.Position1, factors[0])
	}

	counts, err = resultRepo.OneDigit(ctx, fromDate, "second")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by digit-second: %w", err)
	} else {
		digitCount.Position2 = sumProbDigit(counts, digitCount.Position2, factors[1])
	}

	counts, err = resultRepo.OneDigit(ctx, fromDate, "third")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by digit-third: %w", err)
	} else {
		digitCount.Position3 = sumProbDigit(counts, digitCount.Position3, factors[2])
	}

	counts, err = resultRepo.OneDigit(ctx, fromDate, "fourth")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by digit-fourth: %w", err)
	} else {
		digitCount.Position4 = sumProbDigit(counts, digitCount.Position4, factors[3])
	}

	return err
}

func sumProbDigit(results []*model.DigitCount, probAccumulated []float64, factor float64) []float64 {
	totalResults := 0

	for _, p := range results {
		totalResults += p.Count
	}

	for _, p := range results {
		probAccumulated[p.Digit] += (float64(p.Count) / float64(totalResults)) * factor
	}

	return probAccumulated
}

func twoDigitFrequencies(ctx context.Context, fromDate time.Time, resultRepo repository.ResultRepository, twoDigitCount *model.TwoDigitFrequency) error {
	factors := []float64{1.0, 1.45, 1.45, 1.55, 1.55, 5.0}

	counts, err := resultRepo.TwoDigit(ctx, fromDate, "first", "second")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by twoDigit-first-second: %w", err)
	} else {
		twoDigitCount.FirstSecond = sumProbTwoDigit(counts, twoDigitCount.FirstSecond, factors[0])
	}

	counts, err = resultRepo.TwoDigit(ctx, fromDate, "first", "third")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by twoDigit-first-third: %w", err)
	} else {
		twoDigitCount.FirstThird = sumProbTwoDigit(counts, twoDigitCount.FirstThird, factors[1])
	}

	counts, err = resultRepo.TwoDigit(ctx, fromDate, "first", "fourth")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by twoDigit-first-fourth: %w", err)
	} else {
		twoDigitCount.FirstFourth = sumProbTwoDigit(counts, twoDigitCount.FirstFourth, factors[2])
	}

	counts, err = resultRepo.TwoDigit(ctx, fromDate, "second", "third")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by twoDigit-second-third: %w", err)
	} else {
		twoDigitCount.SecondThird = sumProbTwoDigit(counts, twoDigitCount.SecondThird, factors[3])
	}

	counts, err = resultRepo.TwoDigit(ctx, fromDate, "second", "fourth")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by twoDigit-second-fourth: %w", err)
	} else {
		twoDigitCount.SecondFourth = sumProbTwoDigit(counts, twoDigitCount.SecondFourth, factors[4])
	}

	counts, err = resultRepo.TwoDigit(ctx, fromDate, "third", "fourth")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by twoDigit-third-fourth: %w", err)
	} else {
		twoDigitCount.ThirdFourth = sumProbTwoDigit(counts, twoDigitCount.ThirdFourth, factors[5])
	}

	return err
}

func sumProbTwoDigit(results []*model.TwoDigitCount, probAccumulated []float64, factor float64) []float64 {
	totalResults := 0
	weighting := 10.0 // Este weighting es un placeholder, se ajustará usando computación evolutiva to-do-ia

	for _, p := range results {
		totalResults += p.Count
	}

	for _, p := range results {
		num, _ := strconv.Atoi(fmt.Sprintf("%d%d", p.FirstDigit, p.SecondDigit))
		probAccumulated[num] += (float64(p.Count) / float64(totalResults)) * factor * weighting
	}

	return probAccumulated
}

func threeDigitFrequencies(ctx context.Context, fromDate time.Time, resultRepo repository.ResultRepository, threeDigitCount *model.ThreeDigitFrequency) error {
	factors := []float64{1.0, 1.0, 5.0, 5.0}

	counts, err := resultRepo.ThreeDigit(ctx, fromDate, "first", "second", "third")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by threeDigit-first-second-third: %w", err)
	} else {
		threeDigitCount.FirstSecondThird = sumProbThreeDigit(counts, threeDigitCount.FirstSecondThird, factors[0])
	}

	counts, err = resultRepo.ThreeDigit(ctx, fromDate, "first", "second", "fourth")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by threeDigit-first-second-fourth: %w", err)
	} else {
		threeDigitCount.FirstSecondFourth = sumProbThreeDigit(counts, threeDigitCount.FirstSecondFourth, factors[1])
	}

	counts, err = resultRepo.ThreeDigit(ctx, fromDate, "first", "third", "fourth")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by threeDigit-first-third-fourth: %w", err)
	} else {
		threeDigitCount.FirstThirdFourth = sumProbThreeDigit(counts, threeDigitCount.FirstThirdFourth, factors[2])
	}

	counts, err = resultRepo.ThreeDigit(ctx, fromDate, "second", "third", "fourth")
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by threeDigit-second-third-fourth: %w", err)
	} else {
		threeDigitCount.SecondThirdFourth = sumProbThreeDigit(counts, threeDigitCount.SecondThirdFourth, factors[3])
	}

	return nil
}

func sumProbThreeDigit(results []*model.ThreeDigitCount, probAccumulated []float64, factor float64) []float64 {
	totalResults := 0
	weighting := 100.0 // Este weighting es un placeholder, se ajustará usando computación evolutiva to-do-ia

	for _, p := range results {
		totalResults += p.Count
	}

	for _, p := range results {
		num, _ := strconv.Atoi(fmt.Sprintf("%d%d%d", p.FirstDigit, p.SecondDigit, p.ThirdDigit))
		probAccumulated[num] += (float64(p.Count) / float64(totalResults)) * factor * weighting
	}

	return probAccumulated
}

func fourDigitFrequencies(ctx context.Context, fromDate time.Time, resultRepo repository.ResultRepository, fourDigitCount *model.FourDigitFrequency) error {
	counts, err := resultRepo.FourthDigit(ctx, fromDate)
	if err != nil {
		return fmt.Errorf("failed to get frecuencies by all digits: %w", err)
	} else {
		totalResults := 0
		for _, p := range counts {
			totalResults += p.Count
		}
		for _, p := range counts {
			fourDigitCount.Complete[p.Number] += (float64(p.Count) / float64(totalResults)) * 1000
		}
	}

	return err
}
