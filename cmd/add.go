package cmd

import (
	"budgetpipe/cmd/budget"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Import transactions data into your budget on a given month",
	RunE: func(cmd *cobra.Command, args []string) error {
		month := cmd.Flags().Lookup("month").Value.String()

		csvPath, err := getCSVFilePath()
		if err != nil {
			return err
		}

		xlsxPath, err := getXlSXFilePath()
		if err != nil {
			return err
		}

		mapperPath, err := getMapperFilePath()
		if err != nil {
			return err
		}

		return budget.Run(csvPath, xlsxPath, getSheetName(), mapperPath, month)
	},
}

func init() {
	addCmd.Flags().StringP("month", "m", "", "Month to import")
	_ = addCmd.MarkFlagRequired("month")

	rootCmd.AddCommand(addCmd)
}
