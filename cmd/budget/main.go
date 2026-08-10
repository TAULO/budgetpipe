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

const csvTestPath = "./data/test.csv"
const xlsxTestPath = "./data/test.xlsx"

func main() {
	data, err := os.ReadFile("cmd/budget/mapper.json")
	if err != nil {
		panic(err)
	}
	var mapper Mapper
	err = json.Unmarshal(data, &mapper)
	if err != nil {
		panic(err)
	}

	workbook, err := xlsx.NewBudget(xlsxTestPath, "Sheet1")
	if err != nil {
		panic(err)
	}
	defer workbook.Close()

	categories, err := workbook.GetCategories()
	if err != nil {
		panic(err)
	}

	for _, category := range categories {
		expenses := mapper.Expenses[category]
		incomes := mapper.Income[category]
		allCategories := append(expenses, incomes...)

		amount := csv.GetTotalAmountForCategory(csvTestPath, allCategories)
		if amount < 0 {
			amount *= -1
		}
		var _ = workbook.WriteCellFloat(category, "August", amount)

		fmt.Println(category, incomes)
	}

	workbook.Save()
	workbook.Close()
}
