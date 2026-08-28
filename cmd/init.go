package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new budget",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := getWorkDir()
		if err != nil {
			return err
		}

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating work dir %s: %w", dir, err)
		}

		d, err := os.Open(dir)
		if err != nil {
			return fmt.Errorf("opening work dir: %w", err)
		}
		entries, err := d.Readdirnames(1)
		d.Close()
		if err != nil && err != io.EOF {
			return fmt.Errorf("reading work dir: %w", err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("budget already initialized in %s", dir)
		}

		wb := excelize.NewFile()
		budgetFile, err := getXlSXFile()
		if err := wb.SetSheetName("Sheet1", getSheetname()); err != nil {
			return fmt.Errorf("writing work sheet: %w", err)
		}
		if err := wb.SaveAs(budgetFile); err != nil {
			return fmt.Errorf("creating workbook: %w", err)
		}
		wb.Close()

		mapperFile, err := getMapperFile()
		starterMapper := []byte(`{
  "expenses": {},
  "income": {},
  "ignore": []
}`)
		if err := os.WriteFile(mapperFile, starterMapper, 0o644); err != nil {
			return fmt.Errorf("creating mapper: %w", err)
		}

		fmt.Printf("initialized budget in %s\n", dir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
