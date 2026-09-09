package cmd

import (
	"budgetpipe/internal/csv"
	"budgetpipe/internal/xlsx"
	"budgetpipe/tables"
	"fmt"

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

		reader, err := csv.NewReader(monthFilePath)
		if err != nil {
			return fmt.Errorf("creating csv reader: %w", err)
		}

		budget, err := xlsx.NewBudget(budgetPath, getSheetName())
		if err != nil {
			return fmt.Errorf("creating budget: %w", err)
		}

		budgetCategories, err := budget.TableCategories()

		fmt.Println(reader.Transactions())

		fmt.Println(budget.WriteCellByCategoryAndMonth(tables.Income, budgetCategories.Income[0], month, 2800))
		fmt.Println(budget.WriteCellByCategoryAndMonth(tables.Variable, budgetCategories.Variable[0], month, -2000))
		fmt.Println(budget.WriteCellByCategoryAndMonth(tables.Fixed, budgetCategories.Fixed[0], month, -3000))

		return budget.Save()
	},
}

func init() {
	addCmd.Flags().StringP("month", "m", "", "Month to import")
	_ = addCmd.MarkFlagRequired("month")

	rootCmd.AddCommand(addCmd)
}
