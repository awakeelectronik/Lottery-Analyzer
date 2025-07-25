package model

type DigitFrequency struct {
	Position1 []float64 `json:"position_1"`
	Position2 []float64 `json:"position_2"`
	Position3 []float64 `json:"position_3"`
	Position4 []float64 `json:"position_4"`
}

type TwoDigitFrequency struct {
	FirstSecond  []float64 `json:"first_second"`
	FirstThird   []float64 `json:"first_third"`
	FirstFourth  []float64 `json:"first_fourth"`
	SecondThird  []float64 `json:"second_third"`
	SecondFourth []float64 `json:"second_fourth"`
	ThirdFourth  []float64 `json:"third_fourth"`
}

type ThreeDigitFrequency struct {
	FirstSecondThird  []float64 `json:"first_second_third"`
	FirstSecondFourth []float64 `json:"second_fourth_first"`
	FirstThirdFourth  []float64 `json:"third_fourth_first"`
	SecondThirdFourth []float64 `json:"second_third_fourth"`
}

type FourDigitFrequency struct {
	Complete []float64 `json:"complete"`
}

type AllDigitFrequencies struct {
	Position1 []DigitCount `json:"position_1"`
	Position2 []DigitCount `json:"position_2"`
	Position3 []DigitCount `json:"position_3"`
	Position4 []DigitCount `json:"position_4"`
}

type FrequencyData struct {
	DigitFreq      *AllDigitFrequencies `json:"digit_frequencies"`
	TwoDigitFreq   []TwoDigitCount      `json:"two_digit_frequencies"`
	ThreeDigitFreq []ThreeDigitCount    `json:"three_digit_frequencies"`
	FourDigitFreq  []FourDigitCount     `json:"four_digit_frequencies"`
}

type DatabaseStats struct {
	TotalConnections int `json:"total_connections"`
	OpenConnections  int `json:"open_connections"`
	IdleConnections  int `json:"idle_connections"`
	InUseConnections int `json:"in_use_connections"`
}

type TableInfo struct {
	TableName string `json:"table_name"`
	RowCount  int64  `json:"row_count"`
	DataSize  string `json:"data_size"`
}
