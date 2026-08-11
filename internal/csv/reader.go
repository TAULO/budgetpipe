package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/guid"
)

import "budgetpipe/internal/model"

func TotalForCategories(transactions []model.Transaction, categories []string) int64 {
	var total int64
	for _, t := range transactions {
		if slices.Contains(categories, strings.TrimSpace(t.Category)) {
			total += t.Amount
		}
	}
	return total
}

func Transactions(path string) ([]model.Transaction, error) {
	reader, err := readCSVFile(path)
	var transactions []model.Transaction

	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()

		if err != nil || err == io.EOF {
			break
		}

		amount, err := parseDanishAmount(record[2])

		if err != nil {
			//fmt.Println("Error parsing amount:", err)
			continue
		}

		date, err := time.Parse("02.01.2006", record[0])
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		transaction := model.Transaction{
			ID:       guid.New().String(),
			Date:     date,
			Amount:   amount,
			Category: record[8],
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func parseDanishAmount(value string) (int64, error) {
	value = strings.TrimSpace(value)

	// Remove thousands separator
	value = strings.ReplaceAll(value, ".", "")

	// Convert decimal comma to decimal point
	value = strings.ReplaceAll(value, ",", ".")

	// Parse as float temporarily, then convert to øre
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return int64(math.Round(amount * 100)), nil
}

func readCSVFile(path string) (*csv.Reader, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'

	return reader, err
}
