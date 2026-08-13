package cmd

import (
	"budgetpipe/cmd/budget"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// for testing:
// budget  --csv "data/csv/august.csv" --xlsx "data/test.xlsx" --month august

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Import transactions into your budget",
	RunE: func(cmd *cobra.Command, args []string) error {
		csvPath := cmd.Flags().Lookup("csv").Value.String()
		xlsxPath := cmd.Flags().Lookup("xlsx").Value.String()
		month := cmd.Flags().Lookup("month").Value.String()

		mapperPath := defaultMapperPath()

		return budget.Run(csvPath, xlsxPath, mapperPath, month)
	},
}

func Execute() error {
	return budgetCmd.Execute()
}

func init() {
	budgetCmd.Flags().String("csv", "", "path to the bank data CSV file")
	budgetCmd.Flags().String("xlsx", "", "path to the budget Excel file")
	budgetCmd.Flags().String("month", "", "month to import")

	_ = budgetCmd.MarkFlagRequired("csv")
	_ = budgetCmd.MarkFlagRequired("xlsx")
	_ = budgetCmd.MarkFlagRequired("month")
}

func defaultMapperPath() string {
	path, _ := os.Getwd()

	return filepath.Join(
		path,
		"config",
		"mapper.json",
	)
}
