package main

import (
	"budgetpipe/internal/csv"
	"budgetpipe/internal/xlsx"
	"fmt"
)

const csvTestPath = "./data/test.csv"
const xlsxTestPath = "./data/test.xlsx"

func main() {
	// xlsx (Kategori) -> csv (Kategori)
	categoryMap := map[string]string{
		"Daglivarer":          "Dagligvarer",
		"Byen / Cafe":         "Koncert, biograf og museum",
		"Transport":           "Taxa og offentlig transport",
		"Andet (Overførelse)": "Andet (Overførsel)",
		"Abonnementer":        "Film, musik, apps og software",
		"Internet og Telefon": "Telefon, internet, streaming og TV",
		"El (NRGI)":           "El",
		"A-kasse":             "Fagforening, A-kasse og lønsikring",
		"Fitness":             "Sport og fritidsaktiviteter",
		"Andet (Diverse)":     "Lån og gæld (Andet)",
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
		expense := csv.GetExpensesForCategory(csvTestPath, categoryMap[category])
		workbook.WriteCellFloat(category, "August", expense*-1)
		fmt.Println(category, expense, err)
	}

	workbook.Save()
	workbook.Close()
}
