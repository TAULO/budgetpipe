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
	transactions []model.Transaction
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
	r.transactions = transactions
	return r, nil
}

func (r *Reader) ValidateTransactions() error {
	transactions := r.transactions

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

func (r *Reader) Transactions() []model.Transaction {
	return r.transactions
}

func (r *Reader) TransactionsByCategories(categories []string) []model.Transaction {
	var transactions []model.Transaction

	for _, t := range r.transactions {
		category := strings.TrimSpace(t.Category)

		if slices.Contains(categories, category) {
			transactions = append(transactions, t)
		}
	}

	return transactions
}

func (r *Reader) Total(transactions []model.Transaction) int64 {
	var total int64

	for _, transaction := range transactions {
		total += transaction.Amount
	}

	return total
}

func (r *Reader) Comments(transactions []model.Transaction) []string {
	var comments []string

	for _, transaction := range transactions {
		comments = append(comments, transaction.Comment)
	}

	return comments
}

func (r *Reader) UnmappedComments(transactions []model.Transaction) string {
	var comments []string

	for _, transaction := range transactions {
		comments = append(
			comments,
			fmt.Sprintf(
				"• %s: %.2f kr.",
				transaction.Category,
				formatDanishAmount(transaction.Amount),
			),
		)
	}

	return strings.Join(comments, "\n")
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
			Comment:  record[9],
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

func formatDanishAmount(danishCent int64) float64 {
	return float64(danishCent) / 100
}
