package model

type Result struct {
	ID      int    `json:"id" db:"id"`
	Version int    `json:"version" db:"version"`
	Date    string `json:"date" db:"date"`
	First   int    `json:"first" db:"first"`
	Second  int    `json:"second" db:"second"`
	Third   int    `json:"third" db:"third"`
	Fourth  int    `json:"fourth" db:"fourth"`
	Sign    string `json:"sign" db:"sign"`
}

type DigitCount struct {
	Digit int `json:"digit"`
	Count int `json:"count"`
}

type TwoDigitCount struct {
	Count       int `json:"count"`
	FirstDigit  int `json:"first_digit"`
	SecondDigit int `json:"second_digit"`
}

type ThreeDigitCount struct {
	Count       int `json:"count"`
	FirstDigit  int `json:"first_digit"`
	SecondDigit int `json:"second_digit"`
	ThirdDigit  int `json:"third_digit"`
}

type FourDigitCount struct {
	Number int `json:"number"`
	Count  int `json:"count"`
}
