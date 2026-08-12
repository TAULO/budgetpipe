package budget

import (
	"budgetpipe/internal/csv"
	"budgetpipe/internal/model"
	"budgetpipe/internal/xlsx"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type Mapper struct {
	Expenses map[string][]string `json:"expenses"`
	Income   map[string][]string `json:"income"`
}

const fallbackCategory = "Andet (Diverse)"

func main() {
	const month = "August"
	const csvTestPath = "./data/csv/" + month + ".csv"
	const xlsxTestPath = "./data/test.xlsx"
	const mapperPath = "config/mapper.json"

	err := Run(csvTestPath, xlsxTestPath, mapperPath, month)
	if err != nil {
		fmt.Println(err)
	}
}

func Run(csvPath string, xlsxPath string, mapperPath string, month string) error {
	data, err := os.ReadFile(mapperPath)
	if err != nil {
		return fmt.Errorf("reading mapper: %w", err)
	}

	var mapper Mapper
	err = json.Unmarshal(data, &mapper)
	if err != nil {
		return fmt.Errorf("unmarshaling mapper: %w", err)
	}

	workbook, err := xlsx.NewBudget(xlsxPath, "Sheet1")
	if err != nil {
		return fmt.Errorf("opening workbook: %w", err)
	}
	defer workbook.Close()

	reader, err := csv.NewReader(csvPath)
	if err != nil {
		return fmt.Errorf("reading transactions: %w", err)
	}

	if err := reader.ValidateTransactions(); err != nil {
		return err
	}

	write := func(category string, amount int64) error {
		if err := workbook.WriteCellFloat(category, month, amount); err != nil {
			return fmt.Errorf("writing %s: %w", category, err)
		}
		return nil
	}

	writeWithComment := func(category string, amount int64, comment string) error {
		if err := workbook.WriteCellFloatWithComment(category, month, amount, comment); err != nil {
			return fmt.Errorf("writing %s: %w", category, err)
		}
		return nil
	}

	slog.Info("writing transaction to budget", "path", xlsxPath)

	totals := make(map[string]int64)
	var unmapped []model.Transaction

	for _, transaction := range reader.Transactions() {
		csvCategory := strings.TrimSpace(transaction.Category)

		mapped := false

		for budgetCategory, bankCategories := range mapper.Expenses {
			if slices.Contains(bankCategories, csvCategory) {
				totals[budgetCategory] += -transaction.Amount
				mapped = true
				break
			}
		}

		if mapped {
			continue
		}

		for budgetCategory, bankCategories := range mapper.Income {
			if slices.Contains(bankCategories, csvCategory) {
				totals[budgetCategory] += transaction.Amount
				mapped = true
				break
			}
		}

		if !mapped {
			unmapped = append(unmapped, transaction)
		}
	}

	for category, total := range totals {
		if err := write(category, total); err != nil {
			return err
		}
	}

	unmappedTotal := reader.Total(unmapped)
	unmappedComments := reader.UnmappedComments(unmapped)
	if err := writeWithComment(fallbackCategory, -unmappedTotal, unmappedComments); err != nil {
		return err
	}

	for _, transaction := range unmapped {
		fmt.Println("UNMAPPED:", transaction)
	}

	for category := range mapper.Income {
		if err := write(category, totals[category]); err != nil {
			return err
		}
	}

	for category := range mapper.Expenses {
		if err := write(category, totals[category]); err != nil {
			return err
		}
	}

	for category, amount := range totals {
		slog.Info(
			"writing transaction to budget",
			"category", category,
			"amount", amount,
		)
	}

	_ = workbook.Save()

	return nil
}
