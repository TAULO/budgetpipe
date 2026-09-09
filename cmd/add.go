package cmd

import (
	"budgetpipe/internal/csv"
	"budgetpipe/internal/model"
	"budgetpipe/internal/xlsx"
	"budgetpipe/tables"
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Import transactions data into your budget on a given month",
	RunE: func(cmd *cobra.Command, args []string) error {
		month := cmd.Flags().Lookup("month").Value.String()

		monthFilePath, err := getDataCSVFilePath(month)
		if err != nil {
			return fmt.Errorf("resolving csv data path: %w", err)
		}

		budgetPath, err := getBudgetFilePath()
		if err != nil {
			return fmt.Errorf("resolving xlsx data path: %w", err)
		}

		mapperPath, err := getMapperFilePath()
		if err != nil {
			return fmt.Errorf("resolving csv data path: %w", err)
		}

		reader, err := csv.NewReader(monthFilePath)
		if err != nil {
			return fmt.Errorf("creating csv reader: %w", err)
		}

		budget, err := xlsx.NewBudget(budgetPath, getSheetName())
		if err != nil {
			return fmt.Errorf("creating budget: %w", err)
		}

		store := model.NewMapperStore(mapperPath)
		mapper, err := store.Load()
		if err != nil {
			return fmt.Errorf("loading mapper: %w", err)
		}

		income := make(map[string]int64)
		variable := make(map[string]int64)
		fixed := make(map[string]int64)

		for _, transaction := range reader.Transactions() {
			csvCategory := strings.TrimSpace(transaction.Category)

			for budgetCategory, bankCategories := range mapper.Income.Categories {
				if slices.Contains(bankCategories, csvCategory) {
					income[budgetCategory] += transaction.Amount
					break
				}
			}

			for budgetCategory, bankCategories := range mapper.Expenses.Variable {
				if slices.Contains(bankCategories, csvCategory) {
					variable[budgetCategory] += transaction.Amount
					break
				}
			}

			for budgetCategory, bankCategories := range mapper.Expenses.Fixed {
				if slices.Contains(bankCategories, csvCategory) {
					fixed[budgetCategory] += transaction.Amount
					break
				}
			}
		}

		for category, total := range income {
			fmt.Printf("%s: %d\n", category, total)
			if err := budget.WriteToCellInTable(xlsx.CellInput{
				TableName: tables.Income,
				Category:  category,
				Month:     month,
				Value:     total,
			}); err != nil {
				return fmt.Errorf("writing income: %w", err)
			}
		}

		//
		//unmappedTotal := reader.Total(unmapped)
		//unmappedComments := reader.UnmappedComments(unmapped)
		//budget.WriteToCellInTable()
		//
		//for category := range mapper.Income {
		//	if err := write(category, totals[category], ""); err != nil {
		//		return err
		//	}
		//}
		//
		//for category := range mapper.Expenses {
		//	if err := write(category, totals[category], ""); err != nil {
		//		return err
		//	}
		//}
		//
		//for _, transaction := range unmapped {
		//	fmt.Println("UNMAPPED:", transaction)
		//}

		return budget.Save()
	},
}

func init() {
	addCmd.Flags().StringP("month", "m", "", "Month to import")
	_ = addCmd.MarkFlagRequired("month")

	rootCmd.AddCommand(addCmd)
}
