package model

import "time"

type Transaction struct {
	ID       string
	Date     time.Time
	Amount   int64
	Category string
}
