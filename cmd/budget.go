package cmd

import (
	"budgetpipe/cmd/budget"
	"os"

	"github.com/spf13/cobra"
)

// for testing:
// budget  --csv "data/csv/august.csv" --xlsx "data/test.xlsx" --mapper "config/mapper.json" --month december

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Import transactions into your budget",
	RunE: func(cmd *cobra.Command, args []string) error {
		csvPath := cmd.Flags().Lookup("csv").Value.String()
		xlsxPath := cmd.Flags().Lookup("xlsx").Value.String()
		mapperPath := cmd.Flags().Lookup("mapper").Value.String()
		month := cmd.Flags().Lookup("month").Value.String()

		return budget.Run(csvPath, xlsxPath, mapperPath, month)
	},
}

func Execute() error {
	return budgetCmd.Execute()
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	budgetCmd.Flags().StringP("csv", "", home, "path to the csv file")
	budgetCmd.Flags().StringP("xlsx", "", home, "path to budget xlsx file")
	budgetCmd.Flags().StringP("mapper", "", "cmd/budget/mapper.json", "path to mapper")
	budgetCmd.Flags().StringP("month", "", "", "month to import")
}
