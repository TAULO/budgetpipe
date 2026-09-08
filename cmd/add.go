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

		xlsxFilePath, err := getXlSXFilePath()
		if err != nil {
			return fmt.Errorf("resolving xlsx data path: %w", err)
		}

		reader, err := csv.NewReader(monthFilePath)
		if err != nil {
			return fmt.Errorf("creating csv reader: %w", err)
		}

		budget, err := xlsx.NewBudget(template, getSheetName())
		if err != nil {
			return fmt.Errorf("creating budget: %w", err)
		}

		tableCategories, err := budget.TableCategories(tables.Income)
		if err != nil {
			return fmt.Errorf("getting table categories: %w", err)
		}

		fmt.Println(tableCategories)
		fmt.Println(budget.WriteCellByCategoryAndMonth(tables.Fixed, tableCategories[0], month, "foo"))
		fmt.Println(reader.Transactions())

		return budget.SaveAs(xlsxFilePath)
	},
}

func init() {
	addCmd.Flags().StringP("month", "m", "", "Month to import")
	_ = addCmd.MarkFlagRequired("month")

	rootCmd.AddCommand(addCmd)
}
