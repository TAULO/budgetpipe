package cmd

import (
	"budgetpipe/internal/model"
	"budgetpipe/internal/xlsx"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var months = []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "syncs the budget expenses and incomes with the mapper.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("sync called")

		mapperFilePath, err := getMapperFilePath()
		if err != nil {
			return fmt.Errorf("getting mapper file: %w", err)
		}

		budgetFilePath, err := getXlSXFilePath()
		if err != nil {
			return fmt.Errorf("getting xlsx file: %w", err)
		}

		data, err := os.ReadFile(mapperFilePath)
		if err != nil {
			return fmt.Errorf("reading mapper: %w", err)
		}

		var mapper model.Mapper
		err = json.Unmarshal(data, &mapper)
		if err != nil {
			return fmt.Errorf("unmarshaling mapper: %w", err)
		}

		workbook, err := xlsx.NewBudget(budgetFilePath, getSheetName())
		if err != nil {
			return fmt.Errorf("creating workbook: %w", err)
		}

		for i, month := range months {
			cell := fmt.Sprintf("%c1", 'A'+i+1)
			_ = workbook.Write(cell, month)
		}

		fmt.Println("expenses:")
		addressIndex := 2
		for category := range mapper.Expenses.Fixed {
			cell := fmt.Sprintf("A%d", addressIndex)
			if err := workbook.Write(cell, category); err != nil {
				return fmt.Errorf("writing %s: %w", category, err)
			}
			addressIndex++
		}

		for category := range mapper.Expenses.Variable {
			cell := fmt.Sprintf("A%d", addressIndex)
			if err := workbook.Write(cell, category); err != nil {
				return fmt.Errorf("writing %s: %w", category, err)
			}
			addressIndex++
		}

		_ = workbook.Save()

		fmt.Println("income:")
		for category, _ := range mapper.Income.Categories {
			fmt.Println(category)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
