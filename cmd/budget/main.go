package main

import (
	"budgetpipe/internal/csv"
	"budgetpipe/internal/xlsx"
	"fmt"
)

const csvTestPath = "./data/test.csv"
const xlsxTestPath = "./data/test.xlsx"

func main() {
	m := map[string]string{
		"Dagligvarer":                 "Daglivarer",
		"Taxa og offentlig transport": "Transport",
		"Andet (Overførsel)":          "Andet (Overførelse)",
	}

	workbook, err := xlsx.NewBudget(xlsxTestPath, "Sheet1")
	if err != nil {
		panic(err)
	}

	for k, v := range m {
		address, err := workbook.DateCellAddressByCategory(v, "August")
		expense := csv.GetExpensesForCategory(csvTestPath, k)
		workbook.WriteCellFloat(address, expense*-1)

		fmt.Println(v, address, err)
	}

	workbook.Save()
	workbook.Close()
}
