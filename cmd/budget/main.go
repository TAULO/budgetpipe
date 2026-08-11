package main

import (
	"budgetpipe/internal/csv"
	"budgetpipe/internal/xlsx"
	"encoding/json"
	"fmt"
	"os"
)

type Mapper struct {
	Expenses map[string][]string `json:"expenses"`
	Income   map[string][]string `json:"income"`
}

const month = "august"
const csvTestPath = "./data/csv/" + month + ".csv"

const xlsxTestPath = "./data/test.xlsx"

func main() {
	err := Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Run() error {
	data, err := os.ReadFile("cmd/budget/mapper.json")
	if err != nil {
		return fmt.Errorf("reading mapper: %w", err)
	}

	var mapper Mapper
	err = json.Unmarshal(data, &mapper)
	if err != nil {
		return fmt.Errorf("unmarshaling mapper: %w", err)
	}

	workbook, err := xlsx.NewBudget(xlsxTestPath, "Sheet1")
	if err != nil {
		return fmt.Errorf("opening workbook: %w", err)
	}
	defer workbook.Close()

	transactions, err := csv.Transactions(csvTestPath)
	if err != nil {
		return fmt.Errorf("reading transactions: %w", err)
	}

	if err := csv.ValidateTransactions(transactions); err != nil {
		return err
	}

	write := func(category string, amount int64) error {
		if err := workbook.WriteCellFloat(category, month, formatDanishAmount(amount)); err != nil {
			return fmt.Errorf("writing %s: %w", category, err)
		}
		return nil
	}

	for category, bankCats := range mapper.Expenses {
		expense := csv.TotalForCategories(transactions, bankCats)
		if err := write(category, -expense); err != nil {
			return err
		}
	}

	for category, bankCats := range mapper.Income {
		income := csv.TotalForCategories(transactions, bankCats)
		if err := write(category, income); err != nil {
			return err
		}
	}

	workbook.Save()
	workbook.Close()

	return nil
}

func formatDanishAmount(danishCent int64) float64 {
	return float64(danishCent) / 100
}
