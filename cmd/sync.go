package cmd

import (
	"budgetpipe/internal/model"
	"budgetpipe/internal/xlsx"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync your budget.xlsx with your mapper.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetPath, err := getBudgetFilePath()
		if err != nil {
			return fmt.Errorf("resolving budget path: %w", err)
		}

		budget, err := xlsx.NewBudget(budgetPath, getSheetName())
		if err != nil {
			return fmt.Errorf("creating budget file: %w", err)
		}

		budgetCategories, err := budget.TableCategories()
		if err != nil {
			return fmt.Errorf("reading budget categories: %w", err)
		}

		starter := model.NewMapper()
		starter.AddCategories(budgetCategories)

		starterMapper, err := json.MarshalIndent(starter, "", "  ")
		if err != nil {
			return fmt.Errorf("building starter mapper: %w", err)
		}

		mapperFilePath, err := getMapperFilePath()
		if err != nil {
			return fmt.Errorf("resolving mapper path: %w", err)
		}

		if err := os.WriteFile(mapperFilePath, starterMapper, 0o644); err != nil {
			return fmt.Errorf("creating mapper: %w", err)
		}

		fmt.Printf("Synced budget.xlsx with your mapper.json\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
