package csv

import (
	"encoding/csv"
	"errors"
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

type Reader struct {
	Transactions []model.Transaction
	reader       *csv.Reader
}

// Options TODO: Add to parser in the future
type Options struct {
	categoryIndex int
	amountIndex   int
	commentIndex  int
}

func NewReader(path string) (*Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cr := csv.NewReader(file)
	cr.Comma = ';'

	r := &Reader{reader: cr}
	transactions, err := r.parse()
	if err != nil {
		return nil, err
	}
	r.Transactions = transactions
	return r, nil
}

func (r *Reader) ValidateTransactions() error {
	transactions := r.Transactions

	if len(transactions) == 0 {
		return errors.New("no transactions found")
	}

	first := transactions[0]
	last := transactions[len(transactions)-1]

	if first.Date.Month() != last.Date.Month() {
		return errors.New("transactions includes different months")
	}

	return nil
}

func (r *Reader) TotalForCategories(categories []string) int64 {
	var total int64
	for _, t := range r.Transactions {
		if slices.Contains(categories, strings.TrimSpace(t.Category)) {
			total += t.Amount
		}
	}
	return total
}

func (r *Reader) parse() ([]model.Transaction, error) {
	var transactions []model.Transaction

	for {
		record, err := r.reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("reading CSV record: %w", err)
		}

		amount, err := parseDanishAmount(record[2])

		if err != nil {
			continue
		}

		date, err := time.Parse("02.01.2006", record[0])
		if err != nil {
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
