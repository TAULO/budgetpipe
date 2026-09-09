package cmd

import (
	"budgetpipe/flags"
	"budgetpipe/internal/csv"
	"budgetpipe/internal/model"
	"budgetpipe/internal/xlsx"
	"budgetpipe/months"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Import transactions data into your budget on a given month",
	RunE: func(cmd *cobra.Command, args []string) error {
		monthFlag, err := cmd.Flags().GetString(flags.Month)
		if err != nil {
			return err
		}

		month, err := months.Parse(monthFlag)
		if err != nil {
			return err
		}

		monthFilePath, err := getDataCSVFilePath(month.Key())
		if err != nil {
			return fmt.Errorf("resolving csv data path: %w", err)
		}

		budgetPath, err := getBudgetFilePath()
		if err != nil {
			return fmt.Errorf("resolving xlsx data path: %w", err)
		}

		mapperPath, err := getMapperFilePath()
		if err != nil {
			return fmt.Errorf("resolving mapper path: %w", err)
		}

		reader, err := csv.NewReader(monthFilePath)
		if err != nil {
			return fmt.Errorf("creating csv reader: %w", err)
		}

		if err := reader.ValidateTransactions(); err != nil {
			return fmt.Errorf("validating %s: %w", monthFilePath, err)
		}

		mapper, err := model.NewMapperStore(mapperPath).Load()
		if err != nil {
			return fmt.Errorf("loading mapper: %w", err)
		}

		result, err := mapper.Aggregate(reader.Transactions())
		if err != nil {
			return fmt.Errorf("mapping transactions: %w", err)
		}

		budget, err := xlsx.NewBudget(budgetPath, getSheetName())
		if err != nil {
			return fmt.Errorf("creating budget: %w", err)
		}
		defer func() {
			if err := budget.Close(); err != nil {
				slog.Error("closing budget", "error", err)
			}
		}()

		for _, entry := range result.Entries {
			if err := budget.WriteToCellInTable(xlsx.CellInput{
				TableName: entry.Table,
				Category:  entry.Category,
				Month:     month.Header(),
				Value:     entry.Amount,
				Comment:   entry.Comment,
			}); err != nil {
				return fmt.Errorf("writing %s: %w", entry.Category, err)
			}
		}

		for _, transaction := range result.Unmapped {
			slog.Warn(
				"unmapped transaction",
				"category", transaction.Category,
				"amount", transaction.DisplayAmount,
				"comment", transaction.Comment,
			)
		}

		if err := budget.Save(); err != nil {
			return fmt.Errorf("saving budget: %w", err)
		}

		fmt.Printf(
			"imported %s: %d cells written, %d ignored, %d unmapped\n",
			month, len(result.Entries), len(result.Ignored), len(result.Unmapped),
		)

		return nil
	},
}

func init() {
	addCmd.Flags().StringP(flags.Month, "m", "", "Month to import")
	_ = addCmd.MarkFlagRequired(flags.Month)

	rootCmd.AddCommand(addCmd)
}
