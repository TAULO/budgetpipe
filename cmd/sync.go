package cmd

import (
	"budgetpipe/internal/model"
	"budgetpipe/internal/xlsx"
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync your budget.xlsx with your mapper.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		resetFlag, err := cmd.Flags().GetBool("reset")
		if err != nil {
			return err
		}

		budgetPath, err := getBudgetFilePath()
		if err != nil {
			return fmt.Errorf("resolving budget path: %w", err)
		}

		mapperFilePath, err := getMapperFilePath()
		if err != nil {
			return fmt.Errorf("resolving mapper path: %w", err)
		}

		budget, err := xlsx.NewBudget(budgetPath, getSheetName())
		if err != nil {
			return fmt.Errorf("creating budget file: %w", err)
		}

		cats, err := budget.TableCategories()
		if err != nil {
			return fmt.Errorf("reading budget categories: %w", err)
		}

		store := model.NewMapperStore(mapperFilePath)

		var mapper model.Mapper
		if resetFlag {
			mapper = model.NewMapper()
			fmt.Println("resetting mapper.json to match the budget...")
		} else {
			mapper, err = store.Load()
			if err != nil {
				return fmt.Errorf("loading mapper: %w", err)
			}
		}

		mapper.SyncCategories(cats.Income, cats.Fixed, cats.Variable)

		if err := store.Save(mapper); err != nil {
			return fmt.Errorf("saving mapper: %w", err)
		}

		fmt.Println("Synced budget.xlsx with your mapper.json")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP("reset", "r", false, "reset your mapper.json back to the original")
}
